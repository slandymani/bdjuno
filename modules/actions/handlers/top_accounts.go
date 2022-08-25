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
				WHERE ab.address IS NOT NULL
				GROUP BY ab.address
				ORDER BY ab.loki_balance DESC
				OFFSET $1`

	pagination := *payload.GetPagination()
	var rows []dbtypes.TopAccountRow
	var err error

	// To avoid nil result when pagination params undefined (by default set to 0)
	if pagination.Limit != 0 {
		stmt = stmt + `
		LIMIT $2`

		err = db.Sqlx.Select(&rows, stmt, pagination.Offset, pagination.Limit)
	} else {
		err = db.Sqlx.Select(&rows, stmt, pagination.Offset)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to select top accounts")
	}

	if len(rows) == 0 {
		return nil, nil
	}

	var totalCount []int64

	countStmt := `SELECT COUNT(*) FROM account_balance`

	err = db.Sqlx.Select(&totalCount, countStmt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select total count")
	}

	return TopAccountsResponse{
		Rows:       &rows,
		TotalCount: totalCount[0],
	}, nil
}

type TopAccountsResponse struct {
	Rows       *[]dbtypes.TopAccountRow `json:"accounts"`
	TotalCount int64                    `json:"total_count"`
}
