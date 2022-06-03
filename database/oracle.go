package database

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/lib/pq"

	"github.com/forbole/bdjuno/v3/types"
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

func (db *Db) SaveDataSource(height int64, timestamp string, msg *oracletypes.MsgCreateDataSource) error {
	stmt := `
INSERT INTO data_source (create_block, edit_block, name, description, executable, fee, owner, sender, timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := db.Sql.Exec(
		stmt, height, height,
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
	stmt := `UPDATE data_source SET edit_block = $1`

	if msg.Name != "[do-not-modify]" {
		stmt += `, name = $2`
	}
	if msg.Description != "[do-not-modify]" {
		stmt += `, description = $3`
	}
	if string(msg.Executable) != "[do-not-modify]" {
		stmt += `, executable = $4`
	}
	if msg.Fee != nil {
		stmt += `, fee = $5`
	}
	if msg.Owner != "[do-not-modify]" {
		stmt += `, owner = $6`
	}
	stmt += `WHERE id = $7`

	_, err := db.Sql.Exec(
		stmt, height, msg.Name,
		msg.Description, string(msg.Executable),
		pq.Array(dbtypes.NewDbCoins(msg.Fee)), msg.Owner, msg.DataSourceID,
	)
	if err != nil {
		return fmt.Errorf("error while editing oracle script: %s, %s", err, msg.Name)
	}

	return nil
}

func (db *Db) SaveOracleScript(height int64, timestamp string, msg *oracletypes.MsgCreateOracleScript) error {
	stmt := `
INSERT INTO oracle_script (create_block, edit_block, name, description, schema, source_code_url, owner, sender, timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := db.Sql.Exec(
		stmt, height, height,
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
	stmt := `UPDATE oracle_script SET edit_block = $1`

	if msg.Name != "[do-not-modify]" {
		stmt += `, name = $2`
	}
	if msg.Description != "[do-not-modify]" {
		stmt += `, description = $3`
	}
	if msg.Schema != "[do-not-modify]" {
		stmt += `, schema = $4`
	}
	if msg.SourceCodeURL != "[do-not-modify]" {
		stmt += `, source_code_url = $5`
	}
	if msg.Owner != "[do-not-modify]" {
		stmt += `, owner = $6`
	}
	stmt += `WHERE id = $7`

	_, err := db.Sql.Exec(
		stmt, height, msg.Name,
		msg.Description, msg.Schema,
		msg.SourceCodeURL, msg.Owner, msg.OracleScriptID,
	)
	if err != nil {
		return fmt.Errorf("error while editing oracle script: %s", err)
	}

	return nil
}

func (db *Db) SaveDataRequest(timestamp string, msg *oracletypes.MsgRequestData) error {
	calldata := base64.StdEncoding.EncodeToString(msg.Calldata)
	stmt := `
INSERT INTO request (oracle_script_id, calldata, ask_count, min_count, client_id, fee_limit, prepare_gas, execute_gas, sender, timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := db.Sql.Exec(
		stmt, msg.OracleScriptID, calldata, msg.AskCount,
		msg.MinCount, msg.ClientID, pq.Array(dbtypes.NewDbCoins(msg.FeeLimit)),
		msg.PrepareGas, msg.ExecuteGas, msg.Sender, timestamp,
	)
	if err != nil {
		return fmt.Errorf("error while storing data request: %s", err)
	}

	return nil
}

func (db *Db) IncrementReportsCount(msg *oracletypes.MsgReportData) error {
	stmt := `UPDATE request SET reports_count = reports_count + 1 WHERE id = $1`
	_, err := db.Sql.Exec(stmt, msg.RequestID)
	if err != nil {
		return fmt.Errorf("error while incrementing request reports: %s", err)
	}

	return nil
}
