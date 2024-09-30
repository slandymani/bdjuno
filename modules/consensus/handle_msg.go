package consensus

import (
	"fmt"

	juno "github.com/forbole/juno/v6/types"
)

func (m *Module) HandleMsg(_ int, _ juno.Message, tx *juno.Transaction) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	err := m.db.SetTxSender(tx)
	if err != nil {
		return fmt.Errorf("error while setting tx senders: %s", err)
	}

	return nil
}
