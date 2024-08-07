package handlers

import (
	"cosmossdk.io/errors"
	"github.com/forbole/callisto/v4/database"
	"github.com/forbole/callisto/v4/modules/actions/types"
)

type proposalStatistics struct {
	Status string `db:"status" json:"status"`
	Count  uint64 `db:"count" json:"count"`
}

func GetVoteProposalsStatistics(_ *types.Context, _ *types.Payload, db *database.Db) (interface{}, error) {
	var response []proposalStatistics
	stmt := `SELECT status, COUNT(*) AS count
FROM proposal
WHERE status IS NOT NULL
GROUP BY status`
	err := db.Sqlx.Select(&response, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select voting statistics")
	}

	return response, nil
}
