package nft

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/callisto/v4/database"
	onftsource "github.com/forbole/callisto/v4/modules/onft/source"
	"github.com/forbole/juno/v6/modules"
)

var (
	_ modules.Module             = &Module{}
	_ modules.MessageModule      = &Module{}
	_ modules.AuthzMessageModule = &Module{}
	_ modules.GenesisModule      = &Module{}
)

type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source onftsource.Source
}

// NewModule returns a new Module instance
func NewModule(source onftsource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "nft"
}
