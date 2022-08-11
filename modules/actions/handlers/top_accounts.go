package handlers

import (
	"github.com/forbole/bdjuno/v3/modules/actions/types"
	"github.com/rs/zerolog/log"
)

func TopAccountsHandler(ctx *types.Context, payload *types.Payload) (interface{}, error) {
	log.Debug().Msg("executing top accounts action")

	_ = `SELECT ab.address, ab.loki_balance, ab.mgeo_balance, ab.all_balances, COUNT(t.sender) 
				FROM account_balance ab
				FULL OUTER JOIN transaction t ON ab.address = t.sender
				GROUP BY t.sender`

	return nil, nil
}
