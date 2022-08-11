package source

import banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

type Source interface {
	GetTopAccounts(height int64) ([]banktypes.Balance, error)
}
