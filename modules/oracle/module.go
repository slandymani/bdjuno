package oracle

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v5/modules"

	"github.com/forbole/bdjuno/v4/database"
	oraclesource "github.com/forbole/bdjuno/v4/modules/oracle/source"
)

var (
	_ modules.Module                   = &Module{}
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
