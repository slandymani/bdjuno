package handlers

import (
	"github.com/forbole/bdjuno/v3/database"
	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/forbole/bdjuno/v3/modules/actions/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func TopAccountsHandler(ctx *types.Context, payload *types.Payload, db *database.Db) (interface{}, error) {
	log.Debug().Msg("executing top accounts action")

	stmt := `SELECT ab.address, ab.loki_balance, ab.mgeo_balance, ab.all_balances, COUNT(t.sender) as tx_number
				FROM account_balance ab
				FULL OUTER JOIN transaction t ON ab.address = t.sender
				GROUP BY t.sender`

	var rows []dbtypes.TopAccountRow
	err := db.Sqlx.Select(&rows, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select top accounts")
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows, nil
}
