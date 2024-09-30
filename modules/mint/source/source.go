package source

import (
	"cosmossdk.io/math"
	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Source interface {
	GetInflation(height int64) (math.LegacyDec, error)
	Params(height int64) (minttypes.Params, error)
	GetTreasuryPool(height int64) (sdk.Coins, error)
}
