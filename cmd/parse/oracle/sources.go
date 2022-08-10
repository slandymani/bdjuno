package oracle

import (
	"errors"
	"fmt"
	modulestypes "github.com/forbole/bdjuno/v3/modules/types"

	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/types/config"
	"github.com/spf13/cobra"
)

// sourcesCmd returns a Cobra command that allows to refresh data sources.
func sourcesCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "sources",
		Short: "Refresh the information about data sources taking them from the latest known height",
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
			//db := database.Cast(parseCtx.Database)

			// Build the oracle module
			//oracleModule := oracle.NewModule(sources.OracleSource, db)

			// Get latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting latest block height: %s", err)
			}

			// Get all data sources
			dataSources, err := sources.OracleSource.GetDataSources(height)
			if err != nil {
				return fmt.Errorf("error while getting data sources: %s", err)
			}

			s := fmt.Sprintf("Norm?: %d", int(dataSources[0].ID))
			return errors.New(s)
		},
	}
}
