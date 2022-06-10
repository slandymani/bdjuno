package oracle

import (
	"fmt"
	"github.com/forbole/bdjuno/v3/modules/utils"
	"github.com/forbole/bdjuno/v3/types"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "oracle").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.updateOracleParams)
	}); err != nil {
		return err
	}

	if _, err := scheduler.Every(1).Hour().Do(func() {
		utils.WatchMethod(m.updateDataProvidersPool)
	}); err != nil {
		return err
	}

	return nil
}

// updateOracleParams fetches from the REST APIs the latest value for
// oracle params, and saves it inside the database.
func (m *Module) updateOracleParams() error {
	log.Debug().Str("module", "oracle").Msg("getting oracle params data")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the params
	params, err := m.source.GetParams(height)
	if err != nil {
		return err
	}

	return m.db.SaveOracleParams(types.NewOracleParams(params, height), height)
}

func (m *Module) updateDataProvidersPool() error {
	log.Debug().Str("module", "staking").Msg("updating data providers pool")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	pool, err := m.source.GetDataProvidersPool(height)
	if err != nil {
		if err != nil {
			return fmt.Errorf("error while getting data providers pool: %s", err)
		}
	}

	err = m.db.SaveDataProvidersPool(height, pool)
	if err != nil {
		return fmt.Errorf("error while setting data providers pool: %s", err)
	}

	return nil
}
