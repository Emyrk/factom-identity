package factom_identity

import (
	"fmt"

	"time"

	"github.com/Emyrk/factom-raw"
	"github.com/FactomProject/factomd/common/identity"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

var IdentityRegisterChain, _ = primitives.HexToHash("888888001750ede0eff4b05f0c3f557890b256450cabbb84cada937f9c258327")

// Controller can search the blockchain for an identity,
// and feed entries into the parses to come up with the state of the identity.
type Controller struct {
	Reader factom_raw.Fetcher
	Parser *Parser
}

func NewAPIController(apiLocation string) *Controller {
	f := new(Controller)
	f.Reader = factom_raw.NewAPIReader(apiLocation)
	f.Parser = NewParser()

	return f
}

// IsWorking will check if the api is working (connection ok)
func (a *Controller) IsWorking() bool {
	_, err := a.Reader.FetchDBlockHead()
	return err == nil
}

// FindIdentity can be given an authority chain ID and build the identity state.
func (c *Controller) FindIdentity(authorityChain interfaces.IHash) (*identity.Identity, error) {
	// ** Step 1 **
	// First we need to determine if the identity is registered. We will have to parse the entire
	// register chain (TODO: Optimize this)
	regEntries, err := c.FetchChainEntriesInCreateOrder(IdentityRegisterChain)
	if err != nil {
		return nil, err
	}

	err = c.Parser.ParseEntryList(regEntries)
	if err != nil {
		return nil, err
	}

	// ** Step 2 **
	// Parse the authority chain id, which will give us the management chain ID
	rootEntries, err := c.FetchChainEntriesInCreateOrder(authorityChain)
	if err != nil {
		return nil, err
	}

	err = c.Parser.ParseEntryList(rootEntries)
	if err != nil {
		return nil, err
	}

	// ** Step 3 **
	// Parse the entries contained in the management chain (if exists!)
	id := c.Parser.GetIdentity(authorityChain)
	if id == nil {
		return nil, fmt.Errorf("Identity was not found")
	}

	// The id stops here
	if id.ManagementChainID.IsZero() {
		return id, nil
	}

	manageEntries, err := c.FetchChainEntriesInCreateOrder(id.ManagementChainID)
	if err != nil {
		return nil, err
	}

	err = c.Parser.ParseEntryList(manageEntries)
	if err != nil {
		return nil, err
	}

	// ** Step 4 **
	// Return the correct identity
	return c.Parser.GetIdentity(authorityChain), nil
}

// IdentityEntry is parsable, as it contains all the needed info
type IdentityEntry struct {
	Entry       interfaces.IEBEntry
	Timestamp   interfaces.Timestamp
	BlockHeight uint32
}

// FetchChainEntriesInCreateOrder will retrieve all entries in a chain in created order
func (c *Controller) FetchChainEntriesInCreateOrder(chain interfaces.IHash) ([]IdentityEntry, error) {
	now := time.Now()
	head, err := c.Reader.FetchHeadIndexByChainID(chain)
	if err != nil {
		return nil, err
	}

	// Get Eblocks
	var blocks []interfaces.IEntryBlock
	next := head
	for {
		if next.IsZero() {
			break
		}

		// Get the EBlock, and add to list to parse
		block, err := c.Reader.FetchEBlock(next)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)

		next = block.GetHeader().GetPrevKeyMR()
	}

	var entries []IdentityEntry
	// Walk through eblocks in reverse order to get entries
	for i := len(blocks) - 1; i >= 0; i-- {
		eb := blocks[i]

		height := eb.GetDatabaseHeight()
		// Get the timestamp
		dblock, err := c.Reader.FetchDBlockByHeight(height)
		if err != nil {
			return nil, err
		}
		ts := dblock.GetTimestamp()

		ehashes := eb.GetEntryHashes()
		for _, e := range ehashes {
			if e.IsMinuteMarker() {
				continue
			}
			entry, err := c.Reader.FetchEntry(e)
			if err != nil {
				return nil, err
			}

			entries = append(entries, IdentityEntry{entry, ts, height})
		}
	}

	fmt.Println(time.Since(now).Seconds())

	return entries, nil
}
