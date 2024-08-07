package remote

import (
	"fmt"

	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v6/node/remote"
	"github.com/pkg/errors"

	oraclesource "github.com/forbole/callisto/v4/modules/oracle/source"
)

var (
	_ oraclesource.Source = &Source{}
)

// Source implements oraclesource.Source based on a remote node
type Source struct {
	*remote.Source
	client oracletypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, client oracletypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// GetParams implements oraclesource.Source
func (s *Source) GetParams(height int64) (oracletypes.Params, error) {
	res, err := s.client.Params(
		remote.GetHeightRequestContext(s.Ctx, height),
		&oracletypes.QueryParamsRequest{},
	)
	if err != nil {
		return oracletypes.Params{}, fmt.Errorf("error while getting oracle params: %s", err)
	}

	return res.Params, nil
}

func (s *Source) GetRequestStatus(height, id int64) (oracletypes.Result, error) {
	res, err := s.client.Request(
		remote.GetHeightRequestContext(s.Ctx, height),
		&oracletypes.QueryRequestRequest{
			RequestId: uint64(id),
		},
	)
	if err != nil {
		return oracletypes.Result{}, fmt.Errorf("error while getting oracle params: %s", err)
	}

	return *res.Result, nil
}

func (s *Source) GetRequestsInfo(height int64) ([]oracletypes.RequestResult, error) {
	var requests []oracletypes.RequestResult
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.client.Requests(
			remote.GetHeightRequestContext(s.Ctx, height),
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

func (s *Source) GetDataSourcesInfo(height int64) ([]oracletypes.DataSource, error) {
	var dataSources []oracletypes.DataSource
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.client.DataSources(
			remote.GetHeightRequestContext(s.Ctx, height),
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
	response, err := s.client.DataSource(
		remote.GetHeightRequestContext(s.Ctx, height),
		&oracletypes.QueryDataSourceRequest{
			DataSourceId: uint64(id),
		},
	)
	if err != nil {
		return oracletypes.DataSource{}, fmt.Errorf("error while loading data source: %s", err)
	}

	res := oracletypes.DataSource{
		//ID:          response.DataSource.ID,
		Owner:       response.DataSource.Owner,
		Name:        response.DataSource.Name,
		Description: response.DataSource.Description,
		Filename:    response.DataSource.Filename,
		Fee:         response.DataSource.Fee,
	}

	return res, nil
}

func (s *Source) GetRequestInfo(height, id int64) (oracletypes.RequestResult, error) {
	response, err := s.client.Request(
		remote.GetHeightRequestContext(s.Ctx, height),
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

func (s *Source) GetOracleScriptInfo(height, id int64) (oracletypes.OracleScript, error) {
	res, err := s.client.OracleScript(
		remote.GetHeightRequestContext(s.Ctx, height),
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
	var oracleScripts []oracletypes.OracleScript
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.client.OracleScripts(
			remote.GetHeightRequestContext(s.Ctx, height),
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
