package database

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/forbole/bdjuno/v3/types"
	"github.com/lib/pq"
	"time"
)

// SaveOracleParams allows to store the given params inside the database
func (db *Db) SaveOracleParams(params types.OracleParams, height int64) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling oracle params: %s", err)
	}

	stmt := `
INSERT INTO oracle_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE oracle_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing oracle params: %s", err)
	}

	return nil
}

func (db *Db) SaveDataSource(dataSourceId, height int64, timestamp string, msg *oracletypes.MsgCreateDataSource) error {
	stmt := `
INSERT INTO data_source (id, create_block, edit_block, name, description, executable, fee, owner, sender, timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (id) DO NOTHING`

	_, err := db.Sql.Exec(
		stmt, dataSourceId, height, height,
		msg.Name, msg.Description,
		string(msg.Executable), pq.Array(dbtypes.NewDbCoins(msg.Fee)),
		msg.Owner, msg.Owner, timestamp,
	)
	if err != nil {
		return fmt.Errorf("error while storing data source: %s", err)
	}

	return nil
}

func (db *Db) EditDataSource(height int64, msg *oracletypes.MsgEditDataSource) error {
	stmt := `
UPDATE data_source
	SET name = CASE WHEN $2 = '[do-not-modify]' THEN data_source.name ELSE $2 END,
		description = CASE WHEN $3 = '[do-not-modify]' THEN data_source.description ELSE $3 END,
		executable = CASE WHEN $4 = '[do-not-modify]' THEN data_source.executable ELSE $4 END,
		fee = CASE WHEN $5 = '' THEN data_source.fee ELSE $6 END,
		edit_block = $1 WHERE id = $7`

	_, err := db.Sql.Exec(
		stmt, height, msg.Name,
		msg.Description, string(msg.Executable),
		msg.Fee.String(), pq.Array(dbtypes.NewDbCoins(msg.Fee)),
		msg.DataSourceID,
	)
	if err != nil {
		return fmt.Errorf("error while editing data source: %s, %s", err, sdk.NewCoins())
	}

	return nil
}

func (db *Db) SaveOracleScript(oracleScriptId, height int64, timestamp string, msg *oracletypes.MsgCreateOracleScript) error {
	stmt := `
INSERT INTO oracle_script (id, create_block, edit_block, name, description, schema, source_code_url, owner, sender, timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (id) DO NOTHING`

	_, err := db.Sql.Exec(
		stmt, oracleScriptId, height, height,
		msg.Name, msg.Description,
		msg.Schema, msg.SourceCodeURL,
		msg.Owner, msg.Owner, timestamp,
	)
	if err != nil {
		return fmt.Errorf("error while storing oracle script: %s", err)
	}

	return nil
}

func (db *Db) EditOracleScript(height int64, msg *oracletypes.MsgEditOracleScript) error {
	stmt := `
UPDATE oracle_script
	SET name = CASE WHEN $2 = '[do-not-modify]' THEN oracle_script.name ELSE $2 END,
		description = CASE WHEN $3 = '[do-not-modify]' THEN oracle_script.description ELSE $3 END,
		schema = CASE WHEN $4 = '[do-not-modify]' THEN oracle_script.schema ELSE $4 END,
		source_code_url = CASE WHEN $5 = '[do-not-modify]' THEN oracle_script.source_code_url ELSE $5 END,
		edit_block = $1 WHERE id = $6`

	_, err := db.Sql.Exec(
		stmt, height, msg.Name,
		msg.Description, msg.Schema,
		msg.SourceCodeURL, msg.OracleScriptID,
	)
	if err != nil {
		return fmt.Errorf("error while editing oracle script: %s", err)
	}

	return nil
}

func (db *Db) SaveDataRequest(requestId, height int64, dataSourceIDs []int64, timestamp, txHash string, msg *oracletypes.MsgRequestData) error {
	calldata := base64.StdEncoding.EncodeToString(msg.Calldata)
	stmt := `
INSERT INTO request (id, oracle_script_id, height, calldata, ask_count, min_count, client_id, fee_limit, prepare_gas, execute_gas, sender, tx_hash, timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
ON CONFLICT (id) DO NOTHING`
	_, err := db.Sql.Exec(
		stmt, requestId, msg.OracleScriptID, height, calldata, msg.AskCount,
		msg.MinCount, msg.ClientID, pq.Array(dbtypes.NewDbCoins(msg.FeeLimit)),
		msg.PrepareGas, msg.ExecuteGas, msg.Sender, txHash, timestamp,
	)
	if err != nil {
		return fmt.Errorf("error while storing data request: %s", err)
	}

	err = db.saveRequestDataSources(requestId, dataSourceIDs)
	if err != nil {
		return fmt.Errorf("error while storing request data sources: %s", err)
	}

	return nil
}

func (db *Db) saveRequestDataSources(requestId int64, dataSourceIDs []int64) error {
	query := `INSERT INTO request_data_source (request_id, data_source_id) VALUES`

	var params []interface{}
	for i, sourceId := range dataSourceIDs {
		vi := i * 2 // number of rows in table
		query += fmt.Sprintf("($%d,$%d),", vi+1, vi+2)
		params = append(params, requestId, sourceId)
	}
	query = query[:len(query)-1] // remove trailing ","

	query += `
	ON CONFLICT (request_id, data_source_id) DO NOTHING`

	_, err := db.Sql.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while storing request data sources: %s", err)
	}

	return nil
}

func (db *Db) SetRequestStatus(result oracletypes.Result) error {
	stmt := `
UPDATE request
	SET resolve_timestamp = CASE WHEN $1 = 1 THEN $2 ELSE request.resolve_timestamp END
WHERE id = $3`

	_, err := db.Sql.Exec(stmt, result.ResolveStatus, time.Unix(result.ResolveTime, 0), result.RequestID)
	if err != nil {
		return fmt.Errorf("error while setting request report timestamp: %s", err)
	}

	return nil
}

func (db *Db) GetUnresolvedRequests() ([]dbtypes.RequestStatus, error) {
	stmt := `SELECT id, height, min_count, reports_count FROM request WHERE resolve_timestamp = 'epoch'`

	var requests []dbtypes.RequestStatus
	if err := db.Sqlx.Select(&requests, stmt); err != nil {
		return nil, fmt.Errorf("error while getting unresolved requests: %s", err)
	}

	return requests, nil
}

func (db *Db) SaveDataReport(msg *oracletypes.MsgReportData, txHash string) error {
	stmt := `
INSERT INTO report (validator, oracle_script_id, tx_hash)
VALUES ($1, $2, $3)
ON CONFLICT (tx_hash) DO NOTHING`

	stmtSelect := `SELECT oracle_script_id FROM request WHERE id = $1`
	var scriptID []int
	if err := db.Sqlx.Select(&scriptID, stmtSelect, msg.RequestID); err != nil {
		return fmt.Errorf("error while getting oracle script name: %s", err)
	}

	_, err := db.Sql.Exec(stmt, msg.Validator, scriptID[0], txHash)
	if err != nil {
		return fmt.Errorf("error while saving request report: %s", err)
	}

	return nil
}

func (db *Db) AddReportCount(id oracletypes.RequestID) error {
	stmtIncrement := `UPDATE request SET reports_count = reports_count + 1 WHERE id = $1`
	_, err := db.Sql.Exec(stmtIncrement, id)
	if err != nil {
		return fmt.Errorf("error while incrementing reports count: %s", err)
	}

	return nil
}

func (db *Db) SetRequestsPerDate(timestamp string) error {
	stmt := `
INSERT INTO requests_per_date (date, requests_number)
VALUES ($1, 1) ON CONFLICT (date) DO UPDATE
	SET requests_number = requests_per_date.requests_number + 1`

	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	_, err = db.Sqlx.Exec(stmt, TimeToUTCDate(t))
	if err != nil {
		return fmt.Errorf("error while setting requests per date: %s", err)
	}

	var requestsSum int64
	err = db.Sqlx.Get(&requestsSum, `SELECT SUM(requests_number) FROM requests_per_date`)
	if err != nil {
		return fmt.Errorf("error while getting requests sum: %s", err)
	}

	stmtTotal := `
INSERT INTO total_requests (date, total_requests_number)
VALUES ($1, $2) ON CONFLICT (date) DO UPDATE
	SET total_requests_number = total_requests.total_requests_number + 1`

	_, err = db.Sqlx.Exec(stmtTotal, TimeToUTCDate(t), requestsSum)
	if err != nil {
		return fmt.Errorf("error while setting total requests: %s", err)
	}

	return nil
}

func (db *Db) SaveDataProvidersPool(height int64, pool sdk.Coins) error {
	stmt := `
INSERT INTO data_providers_pool (coins, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET coins = excluded.coins,
        height = excluded.height
WHERE data_providers_pool.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, pq.Array(dbtypes.NewDbCoins(pool)), height)
	if err != nil {
		return fmt.Errorf("error while storing data providers pool: %s", err)
	}

	return nil
}

func (db *Db) GetRequestCount() (int64, error) {
	var count int64
	err := db.Sqlx.Get(&count, `SELECT COUNT(*) FROM request`)
	return count, err
}

func (db *Db) GetReportCount() (int64, error) {
	var count int64
	err := db.Sqlx.Get(&count, `SELECT COUNT(*) FROM report`)
	return count, err
}
