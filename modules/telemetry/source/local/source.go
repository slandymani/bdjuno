package local

import (
	"fmt"
	telemetrytypes "github.com/ODIN-PROTOCOL/odin-core/x/telemetry/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	telemetrysource "github.com/forbole/bdjuno/v3/modules/telemetry/source"
	"github.com/forbole/juno/v3/node/local"
)

var (
	_ telemetrysource.Source = &Source{}
)

// Source implements slashingsource.Source using a local node
type Source struct {
	*local.Source
	querier telemetrytypes.QueryServer
}

// NewSource implements a new Source instance
func NewSource(source *local.Source, querier telemetrytypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

func (s Source) GetTopAccounts(height int64) ([]banktypes.Balance, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var balances []banktypes.Balance

	pagination := &query.PageRequest{
		Offset: 0,
		Limit:  1000,
	}
	stop := false

	for !stop {
		res, err := s.querier.TopBalances(sdk.WrapSDKContext(ctx), &telemetrytypes.QueryTopBalancesRequest{
			Denom:      "loki",
			Pagination: pagination,
		})
		if err != nil {
			return nil, err
		}

		stop = len(res.Balances) == 0
		pagination.Offset += pagination.Limit
		balances = append(balances, res.Balances...)
	}

	return balances, nil
}
