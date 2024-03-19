package oracle

import (
	"fmt"

	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	juno "github.com/forbole/juno/v5/types"
	"github.com/rs/zerolog/log"
)

func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, vals *tmctypes.ResultValidators,
) error {
	err := m.updateRequestsResolveTime(block.Block.Height)
	if err != nil {
		return fmt.Errorf("error while updating requests resolve time: %s", err)
	}

	return nil
}

func (m *Module) updateRequestsResolveTime(height int64) error {
	log.Debug().Str("module", "oracle").Int64("height", height).
		Msg("updating requests resolve time")

	requests, err := m.db.GetUnresolvedRequests()
	if err != nil {
		return fmt.Errorf("error while getting unresolved requests: %s", err)
	}

	for _, request := range requests {
		if request.ReportsCount >= request.MinCount || height-request.Height > 100 {
			res, err := m.source.GetRequestStatus(height, request.Id)
			if err != nil {
				return fmt.Errorf("error while getting request result: %s", err)
			}

			err = m.db.SetRequestStatus(res)
			if err != nil {
				return fmt.Errorf("error while setting request resolve time: %s", err)
			}
		}
	}
	return nil
}
