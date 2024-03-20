package oracle

import (
	"strconv"

	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v5/types"
	"github.com/pkg/errors"
)

func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *oracletypes.MsgCreateDataSource:
		dataSourceID, err := GetIDValueFromEvents(tx.Events, oracletypes.EventTypeCreateDataSource, oracletypes.AttributeKeyID)
		if err != nil {
			return errors.Wrap(err, "error while parsing data source id")
		}
		return m.handleMsgCreateDataSource(dataSourceID, tx.Height, tx.Timestamp, cosmosMsg)

	case *oracletypes.MsgEditDataSource:
		return m.handleMsgEditDataSource(tx.Height, cosmosMsg)

	case *oracletypes.MsgCreateOracleScript:
		oracleScriptID, err := GetIDValueFromEvents(tx.Events, oracletypes.EventTypeCreateOracleScript, oracletypes.AttributeKeyID)
		if err != nil {
			return errors.Wrap(err, "error while parsing oracle script id")
		}
		return m.handleMsgCreateOracleScript(oracleScriptID, tx.Height, tx.Timestamp, cosmosMsg)

	case *oracletypes.MsgEditOracleScript:
		return m.handleMsgEditOracleScript(tx.Height, tx.Timestamp, cosmosMsg)

	case *oracletypes.MsgRequestData:
		requestID, err := GetIDValueFromEvents(tx.Events, oracletypes.EventTypeRequest, oracletypes.AttributeKeyID)
		if err != nil {
			return errors.Wrap(err, "error while parsing request id")
		}

		dataSources := GetValueFromEvents(tx.Events, oracletypes.EventTypeRawRequest, oracletypes.AttributeKeyDataSourceID)

		if len(dataSources) == 0 {
			return errors.New("Cannot get request data sources")
		}

		dataSourceIds := make([]int64, len(dataSources))
		for i, v := range dataSources {
			dataSourceID, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return errors.Wrap(err, "error while parsing data source id")
			}
			dataSourceIds[i] = dataSourceID
		}
		return m.handleMsgRequestData(requestID, tx.Height, dataSourceIds, tx.Timestamp, tx.TxHash, cosmosMsg)

	case *oracletypes.MsgReportData:
		reportID, err := GetIDValueFromEvents(tx.Events, oracletypes.EventTypeReport, oracletypes.AttributeKeyID)
		if err != nil {
			return errors.Wrap(err, "error while parsing report id")
		}
		return m.handleMsgReportData(cosmosMsg, tx.TxHash, tx.Height, reportID, tx.Timestamp)
	}

	return nil
}

func (m *Module) handleMsgCreateDataSource(dataSourceID, height int64, timestamp string, msg *oracletypes.MsgCreateDataSource) error {
	err := m.db.SaveDataSource(dataSourceID, height, timestamp, msg)
	if err != nil {
		return errors.Wrap(err, "error while saving data source from MsgCreateDataSource")
	}

	return nil
}

func (m *Module) handleMsgEditDataSource(height int64, msg *oracletypes.MsgEditDataSource) error {
	err := m.db.EditDataSource(height, msg)
	if err != nil {
		return errors.Wrap(err, "error while editing data source from MsgEditDataSource")
	}

	return nil
}

func (m *Module) handleMsgCreateOracleScript(oracleScriptID, height int64, timestamp string, msg *oracletypes.MsgCreateOracleScript) error {
	err := m.db.SaveOracleScript(oracleScriptID, height, timestamp, msg)
	if err != nil {
		return errors.Wrap(err, "error while saving oracle script from MsgCreateOracleScript")
	}

	return nil
}

func (m *Module) handleMsgEditOracleScript(height int64, timestamp string, msg *oracletypes.MsgEditOracleScript) error {
	err := m.db.EditOracleScript(height, msg)
	if err != nil {
		return errors.Wrap(err, "error while editing oracle script from MsgEditOracleScript")
	}

	return nil
}

func (m *Module) handleMsgRequestData(requestID, height int64, dataSourceIDs []int64, timestamp, txHash string, msg *oracletypes.MsgRequestData) error {
	countBeforeSaving, err := m.db.GetRequestCount()
	if err != nil {
		return errors.Wrap(err, "error while saving data request from MsgRequestData")
	}

	err = m.db.SaveDataRequest(requestID, height, dataSourceIDs, timestamp, txHash, msg)
	if err != nil {
		return errors.Wrap(err, "error while saving data request from MsgRequestData")
	}

	countAfterSaving, err := m.db.GetRequestCount()
	if err != nil {
		return errors.Wrap(err, "error while saving data request from MsgRequestData")
	}

	if countAfterSaving > countBeforeSaving {
		err := m.db.SetRequestsPerDate(timestamp)
		if err != nil {
			return errors.Wrap(err, "error while setting requests per date")
		}
	}

	return nil
}

func (m *Module) handleMsgReportData(msg *oracletypes.MsgReportData, txHash string, height, reportID int64, timestamp string) error {
	countBeforeSaving, err := m.db.GetReportCount()
	if err != nil {
		return errors.Wrap(err, "error while saving data request from MsgReportData")
	}

	err = m.db.SaveDataReport(msg, txHash, reportID)
	if err != nil {
		return errors.Wrap(err, "error while saving data report from MsgReportData")
	}

	countAfterSaving, err := m.db.GetReportCount()
	if err != nil {
		return errors.Wrap(err, "error while saving data request from MsgReportData")
	}

	if countAfterSaving > countBeforeSaving {
		err = m.db.AddReportCount(msg.RequestID)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetValueFromEvents(events []abcitypes.Event, eventType, key string) []string {
	res := make([]string, 0)

	for _, event := range events {
		if event.Type == eventType {
			for _, attribute := range event.Attributes {
				if attribute.Key == key {
					res = append(res, attribute.Value)
				}
			}
		}
	}

	return res
}

func GetIDValueFromEvents(events []abcitypes.Event, eventType, key string) (int64, error) {
	idValue := GetValueFromEvents(events, eventType, key)

	if len(idValue) == 0 {
		return 0, errors.New("Id value that matches given key not found")
	}

	id, err := strconv.ParseInt(idValue[0], 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "Error while parsing id value")
	}

	return id, nil
}
