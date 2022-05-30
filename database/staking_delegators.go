package database

import (
	"fmt"
	"github.com/forbole/bdjuno/v3/types"
)

func (db *Db) SaveStakingDelegator(data types.Delegator) error {
	stmt := `
INSERT INTO delegator (address, delegations, delegations_percentage, height)
VALUES ($1, $2, $3, $4)
ON CONFLICT (address) DO UPDATE
	SET delegations = excluded.delegations,
    delegations_percentage = excluded.delegations_percentage,
    height = excluded.height
WHERE delegator.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, data.Address, data.Delegations, data.DelegationsPercentage, data.Height)
	if err != nil {
		return fmt.Errorf("error while storing delegator stake: %s", err)
	}

	return nil
}
