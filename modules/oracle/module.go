package oracle

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v6/modules"

	"github.com/forbole/callisto/v4/database"
	oraclesource "github.com/forbole/callisto/v4/modules/oracle/source"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.BlockModule              = &Module{}
	_ modules.MessageModule            = &Module{}
	_ modules.AuthzMessageModule       = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent database/oracle module
type Module struct {
	db     *database.Db
	source oraclesource.Source
	cdc    codec.Codec
}

// NewModule returns a new Module instance
func NewModule(source oraclesource.Source, db *database.Db, cdc codec.Codec) *Module {
	return &Module{
		db:     db,
		source: source,
		cdc:    cdc,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "oracle"
}
