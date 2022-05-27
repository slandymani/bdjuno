package database

import (
	"fmt"
	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/forbole/bdjuno/v3/types"
	"github.com/lib/pq"
)

func (db *Db) SaveStakingDelegator(data types.Delegator) error {
	stmt := `
INSERT INTO delegator (address, delegations, height)
VALUES ($1, $2, $3)
ON CONFLICT (address) DO UPDATE
	SET delegations = excluded.delegations,
    height = excluded.height
WHERE delegator.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, data.Address, pq.Array(dbtypes.NewDbCoins(data.Delegations)), data.Height)
	if err != nil {
		return fmt.Errorf("error while storing delegator stake: %s", err)
	}

	return nil
}
