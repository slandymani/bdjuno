package remote

import (
	"fmt"

	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v5/node/local"
	"github.com/forbole/juno/v5/node/remote"
	"github.com/pkg/errors"

	oraclesource "github.com/forbole/bdjuno/v4/modules/oracle/source"
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
		&oracletypes.QueryRequestRequest{RequestId: uint64(id)},
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

func (s *Source) GetDataSourcesInfo(height int64) ([]oracletypes.DataSource, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var dataSources []oracletypes.DataSource
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.client.DataSources(
			sdk.WrapSDKContext(ctx),
			&oracletypes.QueryDataSourcesRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 data sources at time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while loading data sources: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		dataSources = append(dataSources, res.DataSources...)
	}

	return dataSources, nil
}

func (s *Source) GetDataSourceInfo(height, id int64) (oracletypes.DataSource, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return oracletypes.DataSource{}, fmt.Errorf("error while loading height: %s", err)
	}

	response, err := s.client.DataSource(
		sdk.WrapSDKContext(ctx),
		&oracletypes.QueryDataSourceRequest{
			DataSourceId: uint64(id),
		},
	)
	if err != nil {
		return oracletypes.DataSource{}, fmt.Errorf("error while loading data source: %s", err)
	}

	return *response.DataSource, nil
}

func (s *Source) GetRequestInfo(height, id int64) (oracletypes.RequestResult, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return oracletypes.RequestResult{}, fmt.Errorf("error while loading height: %s", err)
	}

	response, err := s.client.Request(
		sdk.WrapSDKContext(ctx),
		&oracletypes.QueryRequestRequest{
			RequestId: uint64(id),
		},
	)
	if err != nil {
		return oracletypes.RequestResult{}, errors.Wrap(err, "error while loading request")
	}

	res := oracletypes.RequestResult{
		Request: response.Request,
		Result:  response.Result,
		Reports: response.Reports,
	}

	return res, nil
}

func (s *Source) GetRequestsInfo(height int64) ([]oracletypes.RequestResult, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var requests []oracletypes.RequestResult
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.client.Requests(
			sdk.WrapSDKContext(ctx),
			&oracletypes.QueryRequestsRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 requests at time
				},
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "error while loading requests")
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		requests = append(requests, res.Requests...)
	}

	return requests, nil
}

func (s *Source) GetOracleScriptInfo(height, id int64) (oracletypes.OracleScript, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return oracletypes.OracleScript{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.client.OracleScript(
		sdk.WrapSDKContext(ctx),
		&oracletypes.QueryOracleScriptRequest{
			OracleScriptId: uint64(id),
		},
	)
	if err != nil {
		return oracletypes.OracleScript{}, fmt.Errorf("error while getting oracle script result: %s", err)
	}

	return *res.OracleScript, nil
}

func (s *Source) GetOracleScriptsInfo(height int64) ([]oracletypes.OracleScript, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var oracleScripts []oracletypes.OracleScript
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.client.OracleScripts(
			sdk.WrapSDKContext(ctx),
			&oracletypes.QueryOracleScriptsRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 oracle scripts at time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while loading oracle scripts: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		oracleScripts = append(oracleScripts, res.OracleScripts...)
	}

	return oracleScripts, nil
}
