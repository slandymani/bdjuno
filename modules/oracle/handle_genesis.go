package oracle

import (
	"encoding/json"
	"fmt"
	"time"

	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v4/types"
	"github.com/rs/zerolog/log"
)

func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "oracle").Msg("parsing genesis")
	genTime := doc.GenesisTime.Format(time.DateTime)

	// Read the genesis state
	var genState oracletypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[oracletypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading oracle genesis data: %s", err)
	}

	for i, dataSource := range genState.DataSources {
		err = m.db.SaveDataSource(int64(i+1), doc.InitialHeight, genTime, oracletypes.NewMsgCreateDataSource(
			dataSource.Name,
			dataSource.Description,
			[]byte{},
			dataSource.Fee,
			sdk.AccAddress(dataSource.Treasury),
			sdk.MustAccAddressFromBech32(dataSource.Owner),
			sdk.MustAccAddressFromBech32(dataSource.Owner),
		))
		if err != nil {
			return fmt.Errorf("failed to save data source: %s", err)
		}
	}

	for i, oracleScript := range genState.OracleScripts {
		err = m.db.SaveOracleScript(int64(i+1), doc.InitialHeight, genTime, oracletypes.NewMsgCreateOracleScript(
			oracleScript.Name,
			oracleScript.Description,
			oracleScript.Schema,
			oracleScript.SourceCodeURL,
			[]byte{},
			sdk.MustAccAddressFromBech32(oracleScript.Owner),
			sdk.MustAccAddressFromBech32(oracleScript.Owner),
		))
		if err != nil {
			return fmt.Errorf("failed to save oracle script: %s", err)
		}
	}

	err = m.db.SaveOracleParams(types.OracleParams{
		Params: genState.Params,
		Height: doc.InitialHeight,
	}, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("failed to save oracle params: %s", err)
	}

	err = m.db.SaveDataProvidersPool(doc.InitialHeight, genState.OraclePool.DataProvidersPool)
	if err != nil {
		return fmt.Errorf("failed to save data provides pool: %s", err)
	}

	return nil
}
