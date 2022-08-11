package oracle

import (
	"fmt"
	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/oracle"
	modulestypes "github.com/forbole/bdjuno/v3/modules/types"
	"github.com/pkg/errors"
	"strconv"

	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/types/config"
	"github.com/spf13/cobra"
)

// dataSourceCmd returns a Cobra command that allows to refresh data sources by id.
func dataSourceCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "data-source [id]",
		Args:  cobra.ExactValidArgs(1),
		Short: "Refresh the information about selected data source taking it from the latest known height",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse id from args
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

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

			// Get selected data source
			dataSource, err := sources.OracleSource.GetDataSourceInfo(height, int64(id))
			if err != nil {
				return errors.Wrap(err, "error while getting data source")
			}

			// Refresh data source
			err = oracleModule.RefreshDataSourceInfo(height, dataSource)
			if err != nil {
				return errors.Wrap(err, "error while refreshing data source")
			}

			// Everything is ok
			return nil
		},
	}
}
