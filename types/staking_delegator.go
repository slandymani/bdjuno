package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Delegator struct {
	Address     string
	Delegations sdk.Coins
	Height      int64
}

func NewDelegator(address string, stake sdk.Coins, height int64) Delegator {
	return Delegator{
		Address:     address,
		Delegations: stake,
		Height:      height,
	}
}
