package oracle

import (
	"fmt"
	juno "github.com/forbole/juno/v3/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
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

	ids, err := m.db.GetUnresolvedRequests()
	if err != nil {
		return fmt.Errorf("error while getting unresolved requests: %s", err)
	}

	for _, id := range ids {
		res, err := m.source.GetRequestStatus(height, id)
		if err != nil {
			return fmt.Errorf("error while getting request result: %s", err)
		}

		err = m.db.SetRequestStatus(res)
		if err != nil {
			return fmt.Errorf("error while setting request resolve time: %s", err)
		}
	}
	return nil
}
