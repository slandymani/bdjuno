package bank

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog/log"
)

func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "bank").Msg("parsing genesis")

	// Read the genesis state
	var genState banktypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[banktypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading oracle genesis data: %s", err)
	}

	err = m.db.SaveAccountBalances(doc.InitialHeight, genState.Balances)
	if err != nil {
		return fmt.Errorf("error while saving account balances: %s", err)
	}

	return nil
}
