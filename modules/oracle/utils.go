package oracle

import (
	"encoding/hex"
	"github.com/forbole/juno/v3/parser"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
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
