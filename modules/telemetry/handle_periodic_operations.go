package telemetry

import (
	"fmt"
	"github.com/forbole/bdjuno/v3/modules/utils"
	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "telemetry").Msg("setting up periodic tasks")

	// Update accounts balances every 15 mins
	if _, err := scheduler.Every(15).Minutes().Do(func() {
		utils.WatchMethod(m.updateAccountsBalances)
	}); err != nil {
		return fmt.Errorf("error while setting up pricefeed period operations: %s", err)
	}

	return nil
}

func (m *Module) updateAccountsBalances() error {
	log.Debug().Str("module", "telemetry").Msg("updating accounts balances")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	balances, err := m.source.GetTopAccounts(height)
	if err != nil {
		return err
	}

	err = m.db.SaveAccountBalances(height, balances)
	if err != nil {
		return errors.Wrap(err, "failed to save account balances")
	}

	return nil
}
