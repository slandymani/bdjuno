package bank

import (
	"fmt"
	juno "github.com/forbole/juno/v3/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, txs []*juno.Tx, vals *tmctypes.ResultValidators,
) error {
	if len(txs) > 0 {
		err := m.updateBalances(block.Block.Height)
		if err != nil {
			return fmt.Errorf("error while updating balances: %s", err)
		}
	}

	return nil
}
