package remote

import (
	"fmt"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v3/node/remote"
	"github.com/pkg/errors"

	oraclesource "github.com/forbole/bdjuno/v3/modules/oracle/source"
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
		&oracletypes.QueryRequestRequest{RequestId: id},
	)
	if err != nil {
		return oracletypes.Result{}, fmt.Errorf("error while getting oracle params: %s", err)
	}

	return *res.Result, nil
}

func (s *Source) GetDataProvidersPool(height int64) (sdk.Coins, error) {
	res, err := s.client.DataProvidersPool(remote.GetHeightRequestContext(s.Ctx, height), &oracletypes.QueryDataProvidersPoolRequest{})
	if err != nil {
		return oracletypes.QueryDataProvidersPoolResponse{}.Pool, err
	}

	return res.Pool, nil
}

func (s *Source) GetRequests(height int64) ([]oracletypes.RequestResult, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	reqParams := oracletypes.QueryRequestsRequest{Pagination: &query.PageRequest{Limit: 100}}
	res, err := s.client.Requests(
		ctx,
		&reqParams,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error while loading requests")
	}

	return res.Requests, nil
}

func (s Source) GetDataSources(height int64) ([]oracletypes.DataSource, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)
	res, err := s.client.DataSources(
		ctx,
		&oracletypes.QueryDataSourcesRequest{},
	)
	if err != nil {
		return nil, fmt.Errorf("error while loading data sources: %s", err)
	}

	return res.DataSources, nil
}

//TODO:REMOVE---------------
//func (s *Source) GetOracleScriptByRequestId(height, id int64) (oracletypes.OracleScript, error) {
//	req, err := s.GetRequestStatus(height, id)
//	if err != nil {
//		return oracletypes.OracleScript{}, fmt.Errorf("error while getting request result: %s", err)
//	}
//
//	res, err := s.client.OracleScript(
//		remote.GetHeightRequestContext(s.Ctx, height),
//		&oracletypes.QueryOracleScriptRequest{OracleScriptId: int64(req.OracleScriptID)},
//	)
//	if err != nil {
//		return oracletypes.OracleScript{}, fmt.Errorf("error while getting oracle script result: %s", err)
//	}
//
//	return *res.OracleScript, nil
//}
