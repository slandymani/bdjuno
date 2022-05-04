package source

import (
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

type Source interface {
	GetParams(height int64) (oracletypes.Params, error)
}
