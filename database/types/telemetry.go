package types

type TopAccountRow struct {
	Address     string  `db:"address"`
	LokiBalance int64   `db:"loki_balance"`
	MGeoBalance int64   `db:"mgeo_balance"`
	AllBalances DbCoins `db:"all_balances"`
	TxsNumber   int64   `db:"tx_number"`
}
