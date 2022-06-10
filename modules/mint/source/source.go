package source

import (
	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Source interface {
	GetInflation(height int64) (sdk.Dec, error)
	Params(height int64) (minttypes.Params, error)
	GetTreasuryPool(height int64) (sdk.Coins, error)
}
