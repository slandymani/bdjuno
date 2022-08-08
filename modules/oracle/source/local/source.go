package remote

import (
	"fmt"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v3/node/local"
	"github.com/forbole/juno/v3/node/remote"

	oraclesource "github.com/forbole/bdjuno/v3/modules/oracle/source"
)

var (
	_ oraclesource.Source = &Source{}
)

// Source implements oraclesource.Source based on a remote node
type Source struct {
	*local.Source
	client oracletypes.QueryServer
}

// NewSource returns a new Source instance
func NewSource(source *local.Source, client oracletypes.QueryServer) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// GetParams implements oraclesource.Source
func (s *Source) GetParams(height int64) (oracletypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return oracletypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.client.Params(sdk.WrapSDKContext(ctx), &oracletypes.QueryParamsRequest{})
	if err != nil {
		return oracletypes.Params{}, fmt.Errorf("error while getting params: %s", err)
	}

	return res.Params, nil
}

func (s *Source) GetRequestStatus(height, id int64) (oracletypes.Result, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return oracletypes.Result{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.client.Request(
		remote.GetHeightRequestContext(sdk.WrapSDKContext(ctx), height),
		&oracletypes.QueryRequestRequest{RequestId: id},
	)
	if err != nil {
		return oracletypes.Result{}, fmt.Errorf("error while getting request result: %s", err)
	}

	return *res.Result, nil
}

func (s *Source) GetDataProvidersPool(height int64) (sdk.Coins, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return oracletypes.QueryDataProvidersPoolResponse{}.Pool, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.client.DataProvidersPool(sdk.WrapSDKContext(ctx), &oracletypes.QueryDataProvidersPoolRequest{})
	if err != nil {
		return oracletypes.QueryDataProvidersPoolResponse{}.Pool, err
	}

	return res.Pool, nil
}

func (s *Source) GetOracleScriptByRequestId(height, id int64) (oracletypes.OracleScript, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return oracletypes.OracleScript{}, fmt.Errorf("error while loading height: %s", err)
	}

	req, err := s.GetRequestStatus(height, id)
	if err != nil {
		return oracletypes.OracleScript{}, fmt.Errorf("error while getting request result: %s", err)
	}

	res, err := s.client.OracleScript(
		remote.GetHeightRequestContext(sdk.WrapSDKContext(ctx), height),
		&oracletypes.QueryOracleScriptRequest{OracleScriptId: int64(req.OracleScriptID)},
	)
	if err != nil {
		return oracletypes.OracleScript{}, fmt.Errorf("error while getting oracle script result: %s", err)
	}

	return *res.OracleScript, nil
}
