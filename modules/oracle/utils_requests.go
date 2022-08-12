package oracle

import (
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

//RefreshRequestInfos refreshes the info for the request at the provided height
func (m *Module) RefreshRequestInfos(height int64, request oracletypes.RequestResult) error {

	////Converting request into an appropriate form
	//convertedReq := m.convertRequest(height, request)
	//
	////Getting data sources
	//dataSourceIds := make([]int64, len(request.Request.RawRequests))
	//for i, value := range request.Request.RawRequests {
	//	dataSourceIds[i] = int64(value.DataSourceID)
	//}

	//timestamp := time.Unix(int64(request.Request.RequestTime), 0)

	////Saving request
	//err := m.db.SaveDataRequest(
	//	int64(request.Result.RequestID),
	//	height,
	//	dataSourceIds,
	//	timestamp.UTC().Format("2006-01-02 15:04:05"),
	//	&convertedReq,
	//)
	//if err != nil {
	//	return errors.Wrap(err, "error while saving request")
	//}

	return nil
}

func (m *Module) convertRequest(height int64, request oracletypes.RequestResult) oracletypes.MsgRequestData {
	return oracletypes.MsgRequestData{
		OracleScriptID: request.Result.OracleScriptID,
		Calldata:       request.Result.Calldata,
		AskCount:       request.Result.AskCount,
		MinCount:       request.Result.MinCount,
		ClientID:       request.Result.ClientID,
		Sender:         request.Result.ClientID,
		ExecuteGas:     request.Request.ExecuteGas,
	}
}

//Request
//-----------------------------
//ID							+
//Height						+
//OracleScriptID OracleScriptID +
//Calldata       []byte 		+
//AskCount       uint64 		+
//MinCount       uint64 		+
//ClientID       string 		+
//FeeLimit       types.Coins	?
//PrepareGas     uint64			?
//ExecuteGas     uint64 		+
//Sender         string 		? (ClientId empty)
//Tx_hash						?
//Timestamp						+
//Resolve_timestamp				?
//Reports_count					?
//-----------------------------
//Request - data source		    +
