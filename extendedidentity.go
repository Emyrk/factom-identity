package factom_identity

import (
	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/identity"
	"github.com/FactomProject/factomd/common/primitives"
)

type ExtendedIdentity struct {
	IdentityCore identity.Identity `json:"id_core"`
	Extension    IdentityExtension `json:"id_extension"`
}

// IdentityExtension is the unofficial identity fields
type IdentityExtension struct {
	UserCoinbaseAddress string `json:"user_coinbase_address"`
}

func (i *ExtendedIdentity) PopulateExtension() {
	i.Extension.UserCoinbaseAddress = primitives.ConvertFctAddressToUserStr(i.IdentityCore.CoinbaseAddress)
}

type ExtendedAuthority struct {
	AuthorityCore identity.Authority `json:"auth_core"`
	Extension     AuthorityExtension `json:"auth_extension"`
}

// AuthorityExtension is the unofficial authority fields
type AuthorityExtension struct {
	UserCoinbaseAddress string `json:"user_coinbase_address"`
	HumanStatus         string `json:"readable_status"`
}

func (a *ExtendedAuthority) PopulateExtension() {
	a.Extension.UserCoinbaseAddress = primitives.ConvertFctAddressToUserStr(a.AuthorityCore.CoinbaseAddress)
	a.Extension.HumanStatus = constants.IdentityStatusString(a.AuthorityCore.Status)
}
