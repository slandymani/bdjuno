package oracle

import (
	"encoding/hex"

	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/forbole/juno/v5/parser"
)

func (m *Module) HandleOracleTxs(txs []*tmctypes.ResultTx, parseCtx *parser.Context) error {
	for _, tx := range txs {
		transaction, err := parseCtx.Node.Tx(hex.EncodeToString(tx.Tx.Hash()))
		if err != nil {
			return err
		}

		for index, msg := range transaction.GetMsgs() {
			err = m.HandleMsg(index, msg, transaction)
			if err != nil {
				return err
			}
		}
	}

	// Everything is ok
	return nil
}
