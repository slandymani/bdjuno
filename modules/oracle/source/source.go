package source

import (
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Source interface {
	GetParams(height int64) (oracletypes.Params, error)
	GetRequestStatus(height, id int64) (oracletypes.Result, error)
	GetDataProvidersPool(height int64) (sdk.Coins, error)
	GetRequests(height int64) ([]oracletypes.RequestResult, error)
	GetDataSources(height int64) ([]oracletypes.DataSource, error)
	//TODO:REMOVE---------------
	//GetOracleScriptByRequestId(height, id int64) (oracletypes.OracleScript, error)
}
