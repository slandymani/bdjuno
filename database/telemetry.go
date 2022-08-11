package database

import (
	"fmt"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (db *Db) SaveAccountBalances(height int64, balances []banktypes.Balance) error {
	stmt := `INSERT INTO account_balance (address, loki_balance, mgeo_balance, all_balances, height) VALUES`
	var params []interface{}

	for i, balance := range balances {
		bi := i * 5
		stmt += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", bi+1, bi+2, bi+3, bi+4, bi+5)
		params = append(
			params, balance.Address,
			balance.Coins.AmountOf("loki").Int64(),
			balance.Coins.AmountOf("mGeo").Int64(),
			balance.Coins,
			height,
		)
	}

	stmt = stmt[:len(stmt)-1]
	stmt += `
ON CONFLICT (address) DO UPDATE
	SET loki_balance = excluded.loki_balance,
    	mgeo_balance = excluded.mgeo_balance,
		all_balances = excluded.all_balances,
		height = excluded.height
WHERE account_balance.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing accounts balances: %s", err)
	}

	return nil
}
