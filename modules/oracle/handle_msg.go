package oracle

import (
	"fmt"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"
)

func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *oracletypes.MsgCreateDataSource:
		return m.handleMsgCreateDataSource(tx.Height, tx.Timestamp, cosmosMsg)

	case *oracletypes.MsgEditDataSource:
		return m.handleMsgEditDataSource(tx.Height, cosmosMsg)

	case *oracletypes.MsgCreateOracleScript:
		return m.handleMsgCreateOracleScript(tx.Height, tx.Timestamp, cosmosMsg)

	case *oracletypes.MsgEditOracleScript:
		return m.handleMsgEditOracleScript(tx.Height, tx.Timestamp, cosmosMsg)

	case *oracletypes.MsgRequestData:
		return m.handleMsgRequestData(tx.Height, tx.Timestamp, cosmosMsg)

	case *oracletypes.MsgReportData:
		return m.handleMsgReportData(cosmosMsg, tx.TxHash)
	}

	return nil
}

func (m *Module) handleMsgCreateDataSource(height int64, timestamp string, msg *oracletypes.MsgCreateDataSource) error {
	err := m.db.SaveDataSource(height, timestamp, msg)
	if err != nil {
		return fmt.Errorf("error while saving data source from MsgCreateDataSource: %s", err)
	}

	return nil
}

func (m *Module) handleMsgEditDataSource(height int64, msg *oracletypes.MsgEditDataSource) error {
	err := m.db.EditDataSource(height, msg)
	if err != nil {
		return fmt.Errorf("error while editing data source from MsgEditDataSource: %s", err)
	}

	return nil
}

func (m *Module) handleMsgCreateOracleScript(height int64, timestamp string, msg *oracletypes.MsgCreateOracleScript) error {
	err := m.db.SaveOracleScript(height, timestamp, msg)
	if err != nil {
		return fmt.Errorf("error while saving oracle script from MsgCreateOracleScript: %s", err)
	}

	return nil
}

func (m *Module) handleMsgEditOracleScript(height int64, timestamp string, msg *oracletypes.MsgEditOracleScript) error {
	err := m.db.EditOracleScript(height, msg)
	if err != nil {
		return fmt.Errorf("error while editing oracle script from MsgEditOracleScript: %s", err)
	}

	return nil
}

func (m *Module) handleMsgRequestData(height int64, timestamp string, msg *oracletypes.MsgRequestData) error {
	err := m.db.SetRequestsPerDate(timestamp)
	if err != nil {
		return fmt.Errorf("error while setting requests per date: %s", err)
	}

	err = m.db.SaveDataRequest(height, timestamp, msg)
	if err != nil {
		return fmt.Errorf("error while saving data request from MsgRequestData: %s", err)
	}

	return nil
}

func (m *Module) handleMsgReportData(msg *oracletypes.MsgReportData, txHash string) error {
	err := m.db.SaveDataReport(msg, txHash)
	if err != nil {
		return fmt.Errorf("error while saving data report from MsgReportData: %s", err)
	}

	return nil
}
