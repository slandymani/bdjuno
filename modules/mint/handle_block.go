package mint

import (
	"fmt"
	juno "github.com/forbole/juno/v3/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, vals *tmctypes.ResultValidators,
) error {
	err := m.updateTreasuryPool(block.Block.Height)
	if err != nil {
		return fmt.Errorf("error while updating treasury pool: %s", err)
	}

	return nil
}

func (m *Module) updateTreasuryPool(height int64) error {
	log.Debug().Str("module", "mint").Int64("height", height).
		Msg("updating treasury pool")

	pool, err := m.source.GetTreasuryPool(height)
	if err != nil {
		if err != nil {
			return fmt.Errorf("error while getting treasury pool: %s", err)
		}
	}

	err = m.db.SaveTreasuryPool(height, pool)
	if err != nil {
		return fmt.Errorf("error while setting treasury pool: %s", err)
	}

	return nil
}
