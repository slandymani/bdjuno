package oracle

import (
	"fmt"
	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/oracle"
	modulestypes "github.com/forbole/bdjuno/v3/modules/types"
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/types/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// oracleScriptsCmd returns a Cobra command that allows to refresh oracle scripts.
func oracleScriptsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "oracle-scripts",
		Short: "Refresh the information about oracle scripts taking them from the latest known height",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			sources, err := modulestypes.BuildSources(config.Cfg.Node, parseCtx.EncodingConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build the oracle module
			oracleModule := oracle.NewModule(sources.OracleSource, db)

			// Get latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting latest block height: %s", err)
			}

			// Get all oracle scripts
			oracleScripts, err := sources.OracleSource.GetOracleScriptsInfo(height)
			if err != nil {
				return fmt.Errorf("error while getting data sources: %s", err)
			}

			// Refresh each oracle script
			for _, oracleScript := range oracleScripts {
				err = oracleModule.RefreshOracleScriptInfo(height, oracleScript)
				if err != nil {
					return errors.Wrap(err, "error while refreshing oracle scripts")
				}
			}

			// Everything is ok
			return nil
		},
	}
}
