package oracle

import (
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"
	"github.com/pkg/errors"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"strconv"
)

func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *oracletypes.MsgCreateDataSource:
		dataSourceId, err := GetIdValueFromEvents(tx.Events, oracletypes.EventTypeCreateDataSource, oracletypes.AttributeKeyID)
		if err != nil {
			return errors.Wrap(err, "error while parsing data source id")
		}
		return m.handleMsgCreateDataSource(dataSourceId, tx.Height, tx.Timestamp, cosmosMsg)

	case *oracletypes.MsgEditDataSource:
		return m.handleMsgEditDataSource(tx.Height, cosmosMsg)

	case *oracletypes.MsgCreateOracleScript:
		oracleScriptId, err := GetIdValueFromEvents(tx.Events, oracletypes.EventTypeCreateOracleScript, oracletypes.AttributeKeyID)
		if err != nil {
			return errors.Wrap(err, "error while parsing oracle script id")
		}
		return m.handleMsgCreateOracleScript(oracleScriptId, tx.Height, tx.Timestamp, cosmosMsg)

	case *oracletypes.MsgEditOracleScript:
		return m.handleMsgEditOracleScript(tx.Height, tx.Timestamp, cosmosMsg)

	case *oracletypes.MsgRequestData:
		requestId, err := GetIdValueFromEvents(tx.Events, oracletypes.EventTypeRequest, oracletypes.AttributeKeyID)
		if err != nil {
			return errors.Wrap(err, "error while parsing request id")
		}

		dataSources := GetValueFromEvents(tx.Events, oracletypes.EventTypeRawRequest, oracletypes.AttributeKeyDataSourceID)

		if len(dataSources) == 0 {
			return errors.New("Cannot get request data sources")
		}

		dataSourceIds := make([]int64, len(dataSources))
		for i, v := range dataSources {
			dataSourceId, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return errors.Wrap(err, "error while parsing data source id")
			}
			dataSourceIds[i] = dataSourceId
		}
		return m.handleMsgRequestData(requestId, tx.Height, dataSourceIds, tx.Timestamp, cosmosMsg)

	case *oracletypes.MsgReportData:
		return m.handleMsgReportData(cosmosMsg, tx.TxHash, tx.Height, tx.Timestamp)
	}

	return nil
}

func (m *Module) handleMsgCreateDataSource(dataSourceId, height int64, timestamp string, msg *oracletypes.MsgCreateDataSource) error {
	err := m.db.SaveDataSource(dataSourceId, height, timestamp, msg)
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

func (m *Module) handleMsgCreateOracleScript(oracleScriptId, height int64, timestamp string, msg *oracletypes.MsgCreateOracleScript) error {
	err := m.db.SaveOracleScript(oracleScriptId, height, timestamp, msg)
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

func (m *Module) handleMsgRequestData(requestId, height int64, dataSourceIDs []int64, timestamp string, msg *oracletypes.MsgRequestData) error {
	err := m.db.SetRequestsPerDate(timestamp)
	if err != nil {
		return errors.Wrap(err, "error while setting requests per date")
	}

	err = m.db.SaveDataRequest(requestId, height, dataSourceIDs, timestamp, msg)
	if err != nil {
		return errors.Wrap(err, "error while saving data request from MsgRequestData")
	}

	return nil
}

func (m *Module) handleMsgReportData(msg *oracletypes.MsgReportData, txHash string, height int64, timestamp string) error {
	err := m.db.SaveDataReport(msg, txHash)
	if err != nil {
		return errors.Wrap(err, "error while saving data report from MsgReportData")
	}

	return nil
}

func GetValueFromEvents(events []abcitypes.Event, eventType, key string) []string {
	res := make([]string, 0)

	for _, event := range events {
		if event.Type == eventType {
			for _, attribute := range event.Attributes {
				if string(attribute.Key) == key {
					res = append(res, string(attribute.Value))
				}
			}
		}
	}

	return res
}

func GetIdValueFromEvents(events []abcitypes.Event, eventType, key string) (int64, error) {
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
