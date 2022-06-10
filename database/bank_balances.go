package database

import (
	"fmt"
	"github.com/forbole/bdjuno/v3/types"
)

func (db *Db) SaveAccountBalances(balances []types.AccountBalance) error {
	stmt := `INSERT INTO account_balance (address, loki_balance, minigeo_balance, height) VALUES`
	var params []interface{}

	for i, balance := range balances {
		bi := i * 4
		stmt += fmt.Sprintf("($%d,$%d,$%d,$%d),", bi+1, bi+2, bi+3, bi+4)
		params = append(
			params, balance.Address,
			balance.Balance.AmountOf("loki").Int64(),
			balance.Balance.AmountOf("minigeo").Int64(),
			balance.Height,
		)
	}

	stmt = stmt[:len(stmt)-1]
	stmt += `
ON CONFLICT (address) DO UPDATE
	SET loki_balance = excluded.loki_balance,
    	minigeo_balance = excluded.minigeo_balance,
		height = excluded.height
WHERE account_balance.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing accounts balances: %s", err)
	}

	return nil
}
