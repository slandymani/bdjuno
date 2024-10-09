package source

import (
	onfttypes "github.com/ODIN-PROTOCOL/odin-core/x/onft/types"
)

type Source interface {
	ClassOwner(height int64, classID string) (string, error)
	NFTs(height int64, classID, owner string) ([]*onfttypes.NFT, error)
	NFT(height int64, classID, id string) (*onfttypes.NFT, error)
	Class(height int64, classID string) (*onfttypes.Class, error)
	Classes(height int64) ([]*onfttypes.Class, error)
	// 	ClassOwner(context.Context, *onfttypes.QueryClassOwnerRequest) (*onfttypes.QueryClassOwnerResponse, error)
	//	NFTs(context.Context, *onfttypes.QueryNFTsRequest) (*onfttypes.QueryNFTsResponse, error)
	//	NFT(context.Context, *onfttypes.QueryNFTRequest) (*onfttypes.QueryNFTResponse, error)
	//	Class(context.Context, *onfttypes.QueryClassRequest) (*onfttypes.QueryClassResponse, error)
	//	Classes(context.Context, *onfttypes.QueryClassesRequest) (*onfttypes.QueryClassesResponse, error)
}
