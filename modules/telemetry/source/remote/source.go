package remote

import (
	telemetrytypes "github.com/ODIN-PROTOCOL/odin-core/x/telemetry/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	telemetrysource "github.com/forbole/bdjuno/v3/modules/telemetry/source"
	"github.com/forbole/juno/v3/node/remote"
)

var (
	_ telemetrysource.Source = &Source{}
)

// Source implements telemetrysource.Source based on a remote node
type Source struct {
	*remote.Source
	querier telemetrytypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier telemetrytypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

func (s Source) GetTopAccounts(height int64) ([]banktypes.Balance, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var balances []banktypes.Balance

	pagination := &query.PageRequest{
		Offset: 0,
		Limit:  1000,
	}
	stop := false

	for !stop {
		res, err := s.querier.TopBalances(ctx, &telemetrytypes.QueryTopBalancesRequest{
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
