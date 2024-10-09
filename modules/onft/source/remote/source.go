package remote

import (
	"fmt"

	onfttypes "github.com/ODIN-PROTOCOL/odin-core/x/onft/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	onftsource "github.com/forbole/callisto/v4/modules/onft/source"
	"github.com/forbole/juno/v6/node/remote"
)

var (
	_ onftsource.Source = &Source{}
)

// Source implements onftsource.Source based on a remote node
type Source struct {
	*remote.Source
	client onfttypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, client onfttypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

func (s Source) ClassOwner(height int64, classID string) (string, error) {
	res, err := s.client.ClassOwner(
		remote.GetHeightRequestContext(s.Ctx, height),
		&onfttypes.QueryClassOwnerRequest{ClassId: classID},
	)
	if err != nil {
		return "", fmt.Errorf("error while getting nft class owner: %s", err)
	}

	return res.Owner, nil
}

func (s Source) NFTs(height int64, classID, owner string) ([]*onfttypes.NFT, error) {
	var nfts []*onfttypes.NFT
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.client.NFTs(
			remote.GetHeightRequestContext(s.Ctx, height),
			&onfttypes.QueryNFTsRequest{
				ClassId: classID,
				Owner:   owner,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100,
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while loading nfts: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		nfts = append(nfts, res.Nfts...)
	}

	return nfts, nil
}

func (s Source) NFT(height int64, classID, id string) (*onfttypes.NFT, error) {
	response, err := s.client.NFT(
		remote.GetHeightRequestContext(s.Ctx, height),
		&onfttypes.QueryNFTRequest{
			ClassId: classID,
			Id:      id,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error while loading nft: %s", err)
	}

	return response.Nft, nil
}

func (s Source) Class(height int64, classID string) (*onfttypes.Class, error) {
	response, err := s.client.Class(
		remote.GetHeightRequestContext(s.Ctx, height),
		&onfttypes.QueryClassRequest{
			ClassId: classID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error while loading nft class: %s", err)
	}

	return response.Class, nil
}

func (s Source) Classes(height int64) ([]*onfttypes.Class, error) {
	var classes []*onfttypes.Class
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.client.Classes(
			remote.GetHeightRequestContext(s.Ctx, height),
			&onfttypes.QueryClassesRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100,
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while loading nft classes: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		classes = append(classes, res.Classes...)
	}

	return classes, nil
}
