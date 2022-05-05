package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/packet-forward-middleware/router/types"
)

// InitGenesis
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	k.SetParams(ctx, state.Params)
}

// ExportGenesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{Params: k.GetParams(ctx)}
}
