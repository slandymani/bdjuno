package source

import (
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Source interface {
	GetParams(height int64) (oracletypes.Params, error)

	GetRequestInfo(height, id int64) (oracletypes.RequestResult, error)
	GetRequestsInfo(height int64) ([]oracletypes.RequestResult, error)
	GetRequestStatus(height, id int64) (oracletypes.Result, error)

	GetDataSourceInfo(height, id int64) (oracletypes.DataSource, error)
	GetDataSourcesInfo(height int64) ([]oracletypes.DataSource, error)

	GetOracleScriptInfo(height, id int64) (oracletypes.OracleScript, error)
	//GetOracleScriptsInfo(height int64) (oracletypes.OracleScript, error)

	GetDataProvidersPool(height int64) (sdk.Coins, error)
}
