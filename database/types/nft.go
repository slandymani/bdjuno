package types

import codectypes "github.com/cosmos/cosmos-sdk/codec/types"

type NFT struct {
	ClassId    string          `db:"class_id"`
	Id         string          `db:"id"`
	Uri        string          `db:"uri"`
	UriHash    string          `db:"uri_hash"`
	Owner      string          `db:"owner"`
	Data       *codectypes.Any `db:"data"`
	MintTxHash string          `db:"mint_tx_hash"`
}
