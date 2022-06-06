package oracle

import (
	"fmt"
	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	juno "github.com/forbole/juno/v3/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, vals *tmctypes.ResultValidators,
) error {
	unresolvedRequests, err := m.db.GetUnresolvedRequests()
	if err != nil {
		return fmt.Errorf("error while getting unresolved requests: %s", err)
	}

	err = m.updateRequestsResolveTime(block.Block.Height, unresolvedRequests)
	if err != nil {
		return fmt.Errorf("error while updating requests resolve time: %s", err)
	}

	return nil
}

func (m *Module) updateRequestsResolveTime(height int64, ids []dbtypes.UnresolvedRequest) error {
	log.Debug().Str("module", "oracle").Int64("height", height).
		Msg("updating requests resolve time")

	for _, id := range ids {
		res, err := m.source.GetRequestStatus(height, id.UresolvedRequestID)
		if err != nil {
			return fmt.Errorf("error while getting request result from MsgReportData: %s", err)
		}

		err = m.db.SetRequestStatus(res)
		if err != nil {
			return fmt.Errorf("error while setting request resolve time from MsgReportData: %s", err)
		}
	}
	return nil
}
