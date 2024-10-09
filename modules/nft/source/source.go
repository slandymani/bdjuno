package source

import (
	"context"

	"cosmossdk.io/x/nft"
)

type Source interface {
	Owner(context.Context, *nft.QueryOwnerRequest) (*nft.QueryOwnerResponse, error)
	Supply(context.Context, *nft.QuerySupplyRequest) (*nft.QuerySupplyResponse, error)
	NFTs(context.Context, *nft.QueryNFTsRequest) (*nft.QueryNFTsResponse, error)
	NFT(context.Context, *nft.QueryNFTRequest) (*nft.QueryNFTResponse, error)
	Class(context.Context, *nft.QueryClassRequest) (*nft.QueryClassResponse, error)
	Classes(context.Context, *nft.QueryClassesRequest) (*nft.QueryClassesResponse, error)
}
