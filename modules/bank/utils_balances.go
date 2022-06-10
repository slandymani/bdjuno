package bank

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

func (m *Module) RefreshBalances(height int64, addresses []string) error {
	balances, err := m.keeper.GetBalances(addresses, height)
	if err != nil {
		return err
	}

	return m.db.SaveAccountBalances(balances)
}

func (m *Module) updateBalances(height int64) error {
	log.Debug().Str("module", "bank").Int64("height", height).
		Msg("updating modules")

	accounts, err := m.db.GetAccounts()
	if err != nil {
		return fmt.Errorf("error while getting accounts: %s", err)
	}

	err = m.RefreshBalances(height, accounts)
	if err != nil {
		return fmt.Errorf("error while updating balances: %s", err)
	}

	return nil
}
