package factom_identity

import (
	"fmt"

	"encoding/hex"

	"github.com/Emyrk/factom-raw"
	"github.com/FactomProject/factomd/common/identity"
	"github.com/FactomProject/factomd/common/identityEntries"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

var IdentityRegisterChain, _ = primitives.HexToHash("888888001750ede0eff4b05f0c3f557890b256450cabbb84cada937f9c258327")

// Controller can search the blockchain for an identity,
// and feed entries into the parses to come up with the state of the identity.
type Controller struct {
	Reader factom_raw.Fetcher
	Parser *IdentityParser
}

func NewAPIController(apiLocation string) *Controller {
	f := new(Controller)
	f.Reader = factom_raw.NewAPIReader(apiLocation)
	f.Parser = NewIdentityParser()

	return f
}

// IsWorking will check if the api is working (connection ok)
func (a *Controller) IsWorking() bool {
	_, err := a.Reader.FetchDBlockHead()
	return err == nil
}

func (c *Controller) FindAllIdentities() (map[string]*identity.Identity, error) {
	// Find all registered identities
	ids, err := c.parseRegisterChain(true)
	if err != nil {
		return nil, err
	}

	for _, chainID := range ids {
		c.parseIdentityChain(chainID)
	}

	humanMap := make(map[string]*identity.Identity)
	for k, v := range c.Parser.IdentityManager.Identities {
		if v.IdentityChainID == nil || v.IdentityChainID.IsZero() {
			continue
		}
		humanMap[hex.EncodeToString(k[:])] = v
	}
	return humanMap, nil
}

func (c *Controller) parseRegisterChain(getHashes bool) ([]interfaces.IHash, error) {
	regEntries, err := c.FetchChainEntriesInCreateOrder(IdentityRegisterChain)
	if err != nil {
		return nil, err
	}

	var ids []interfaces.IHash
	if getHashes {
		for _, e := range regEntries {
			rfi := new(identityEntries.RegisterFactomIdentityStructure)
			err := rfi.DecodeFromExtIDs(e.Entry.ExternalIDs())
			if err != nil {
				continue
			}

			ids = append(ids, rfi.IdentityChainID)
		}
	}

	err = c.Parser.ParseEntryList(regEntries)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// FindIdentity can be given an authority chain ID and build the identity state.
func (c *Controller) FindIdentity(authorityChain interfaces.IHash) (*identity.Identity, error) {
	// ** Step 1 **
	// First we need to determine if the identity is registered. We will have to parse the entire
	// register chain (TODO: Optimize this)
	_, err := c.parseRegisterChain(false)
	if err != nil {
		return nil, err
	}

	// ** Step 2 **
	// Parse the authority chain id,
	err = c.parseIdentityChain(authorityChain)
	if err != nil {
		return nil, err
	}

	// ** Step 3 **
	// Return the correct identity
	return c.Parser.GetIdentity(authorityChain), nil
}

func (c *Controller) parseIdentityChain(identityChainId interfaces.IHash) error {
	// Parse the root
	rootEntries, err := c.FetchChainEntriesInCreateOrder(identityChainId)
	if err != nil {
		return err
	}

	err = c.Parser.ParseEntryList(rootEntries)
	if err != nil {
		return err
	}

	// Parse the entries contained in the management chain (if exists!)
	id := c.Parser.GetIdentity(identityChainId)
	if id == nil {
		return fmt.Errorf("Identity was not found")
	}

	// The id stops here
	if id.ManagementChainID.IsZero() {
		return nil
	}

	manageEntries, err := c.FetchChainEntriesInCreateOrder(id.ManagementChainID)
	if err != nil {
		return err
	}

	err = c.Parser.ParseEntryList(manageEntries)
	if err != nil {
		return err
	}

	return nil
}

// IdentityEntry is parsable, as it contains all the needed info
type IdentityEntry struct {
	Entry       interfaces.IEBEntry
	Timestamp   interfaces.Timestamp
	BlockHeight uint32
}

// FetchChainEntriesInCreateOrder will retrieve all entries in a chain in created order
func (c *Controller) FetchChainEntriesInCreateOrder(chain interfaces.IHash) ([]IdentityEntry, error) {
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

	return entries, nil
}
