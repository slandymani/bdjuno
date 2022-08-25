package types

type TopAccountRow struct {
	Address     string  `db:"address" json:"address"`
	LokiBalance int64   `db:"loki_balance" json:"loki_balance"`
	MGeoBalance int64   `db:"mgeo_balance" json:"mgeo_balance"`
	AllBalances DbCoins `db:"all_balances" json:"all_balances"`
	TxsNumber   int64   `db:"tx_number" json:"tx_number"`
}
