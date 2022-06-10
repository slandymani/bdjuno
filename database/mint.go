package database

import (
	"encoding/json"
	"fmt"
	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/lib/pq"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/v3/types"
)

// SaveInflation allows to store the inflation for the given block height as well as timestamp
func (db *Db) SaveInflation(inflation sdk.Dec, height int64) error {
	stmt := `
INSERT INTO inflation (value, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET value = excluded.value, 
        height = excluded.height 
WHERE inflation.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, inflation.String(), height)
	if err != nil {
		return fmt.Errorf("error while storing inflation: %s", err)
	}

	return nil
}

// SaveMintParams allows to store the given params inside the database
func (db *Db) SaveMintParams(params *types.MintParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling mint params: %s", err)
	}

	stmt := `
INSERT INTO mint_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE mint_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing mint params: %s", err)
	}

	return nil
}

func (db *Db) SaveTreasuryPool(height int64, pool sdk.Coins) error {
	stmt := `
INSERT INTO treasury_pool (coins, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET coins = excluded.coins,
        height = excluded.height
WHERE treasury_pool.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, pq.Array(dbtypes.NewDbCoins(pool)), height)
	if err != nil {
		return fmt.Errorf("error while storing treasury pool: %s", err)
	}

	return nil
}
