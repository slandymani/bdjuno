package oracle

import (
	"fmt"
	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/oracle"
	modulestypes "github.com/forbole/bdjuno/v3/modules/types"
	"github.com/pkg/errors"

	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/types/config"
	"github.com/spf13/cobra"
)

// dataSourcesCmd returns a Cobra command that allows to refresh data sources.
func dataSourcesCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "data-sources",
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
			db := database.Cast(parseCtx.Database)

			// Build the oracle module
			oracleModule := oracle.NewModule(sources.OracleSource, db)

			// Get latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting latest block height: %s", err)
			}

			// Get all data sources
			dataSources, err := sources.OracleSource.GetDataSourcesInfo(height)
			if err != nil {
				return fmt.Errorf("error while getting data sources: %s", err)
			}

			// Refresh data sources
			for _, dataSource := range dataSources {
				err = oracleModule.RefreshDataSourceInfo(height, dataSource)
				if err != nil {
					return errors.Wrap(err, "error while refreshing data sources")
				}
			}

			// Everything is ok
			return nil
		},
	}
}
