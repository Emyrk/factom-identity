package factom_identity

import (
	"fmt"

	"github.com/FactomProject/factomd/common/identity"
	"github.com/FactomProject/factomd/common/interfaces"
)

// Parser can parse identity related entries or admin blocks. It
// can also be extended to allow for additional entry types (such as naming)
type Parser struct {
	*identity.IdentityManager
}

func NewParser() *Parser {
	p := new(Parser)

	p.IdentityManager = identity.NewIdentityManager()
	return p
}

func (p *Parser) ParseEntryList(list []IdentityEntry) error {
	for _, e := range list {
		err := p.ParseEntry(e.Entry, e.BlockHeight, e.Timestamp, true)
		if err != nil {
			return err
		}
	}

	// Parse the remaining
	err := p.ProcessOldEntries()
	if err != nil {
		return err
	}
	return nil
}

// ParseEntry is mostly handled by the IdentityManager, however it can be extended to support additional parsing options (such as naming)
func (p *Parser) ParseEntry(entry interfaces.IEBEntry, dBlockHeight uint32, dBlockTimestamp interfaces.Timestamp, newEntry bool) error {
	err := p.ProcessIdentityEntry(entry, dBlockHeight, dBlockTimestamp, newEntry)
	if err != nil {
		return err
	}

	//
	if entry.GetChainID().String()[:6] != "888888" {
		return fmt.Errorf("Invalic chainID - expected 888888..., got %v", entry.GetChainID().String())
	}
	if entry.GetHash().String() == "172eb5cb84a49280c9ad0baf13bea779a624def8d10adab80c3d007fe8bce9ec" {
		//First entry, can ignore
		return nil
	}

	// Not always the authority chainID, it can be any chain with '8888', so management, authority, or register chain
	chainID := entry.GetChainID()

	extIDs := entry.ExternalIDs()
	if len(extIDs) < 2 {
		//Invalid Identity Chain Entry
		return fmt.Errorf("Invalid Identity Chain Entry")
	}
	if len(extIDs[0]) == 0 {
		return fmt.Errorf("Invalid Identity Chain Entry")
	}
	if extIDs[0][0] != 0 {
		//We only support version 0
		return fmt.Errorf("Invalid Identity Chain Entry version")
	}

	// This is the entry's name. The ones detailed in the identity spec are covered above, we can support additional
	// types here
	switch string(extIDs[1]) {
	case "TODO":
	}

	var _ = chainID

	return nil
}
