package local

import (
	telemetrytypes "github.com/ODIN-PROTOCOL/odin-core/x/telemetry/types"
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
	return nil, nil
}
