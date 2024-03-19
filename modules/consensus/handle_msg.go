package consensus

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v5/types"
)

func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	err := m.db.SetTxSender(tx)
	if err != nil {
		return fmt.Errorf("error while setting tx senders: %s", err)
	}

	return nil
}
