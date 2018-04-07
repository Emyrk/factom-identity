package factom_identity

import (
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/state"
)

type FakeState struct {
	*state.State
}

func (FakeState) AddIdentityFromChainID(hash interfaces.IHash) error {
	return nil
}
