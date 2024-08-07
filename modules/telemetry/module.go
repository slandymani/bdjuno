package telemetry

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/callisto/v4/database"
	telemetrysource "github.com/forbole/callisto/v4/modules/telemetry/source"
	"github.com/forbole/juno/v6/modules"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent x/slashing module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source telemetrysource.Source
}

// NewModule returns a new Module instance
func NewModule(source telemetrysource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "telemetry-new"
}
