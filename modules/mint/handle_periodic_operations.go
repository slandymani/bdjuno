package mint

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/modules/utils"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "mint").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.UpdateInflation)
	}); err != nil {
		return err
	}

	if _, err := scheduler.Every(1).Hour().Do(func() {
		utils.WatchMethod(m.updateTreasuryPool)
	}); err != nil {
		return err
	}

	return nil
}

// updateInflation fetches from the REST APIs the latest value for the
// inflation, and saves it inside the database.
func (m *Module) UpdateInflation() error {
	log.Debug().
		Str("module", "mint").
		Str("operation", "inflation").
		Msg("getting inflation data")

	block, err := m.db.GetLastBlockHeightAndTimestamp()
	if err != nil {
		return err
	}

	// Get the inflation
	inflation, err := m.source.GetInflation(block.Height)
	if err != nil {
		return err
	}

	return m.db.SaveInflation(inflation, block.Height)
}

func (m *Module) updateTreasuryPool() error {
	log.Debug().Str("module", "mint").Msg("updating treasury pool")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	pool, err := m.source.GetTreasuryPool(height)
	if err != nil {
		if err != nil {
			return fmt.Errorf("error while getting treasury pool: %s", err)
		}
	}

	err = m.db.SaveTreasuryPool(height, pool)
	if err != nil {
		return fmt.Errorf("error while setting treasury pool: %s", err)
	}

	return nil
}
