package oracle

import (
	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/oracle"
	"github.com/pkg/errors"
	"strconv"

	modulestypes "github.com/forbole/bdjuno/v3/modules/types"

	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/types/config"
	"github.com/spf13/cobra"
	//"github.com/forbole/bdjuno/v3/database"
	//"github.com/forbole/bdjuno/v3/modules/oracle"
)

// requestCmd returns a Cobra command that allows to refresh request by id.
func requestCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "request [id]",
		Args:  cobra.ExactValidArgs(1),
		Short: "Refresh the information about selected request taking it from the latest known height",
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
				return errors.Wrap(err, "error while getting latest block height")
			}

			// Get selected request
			reqResponse, err := sources.OracleSource.GetRequestInfo(height, int64(id))
			if err != nil {
				return errors.Wrap(err, "error while getting request")
			}

			//Refresh request
			err = oracleModule.RefreshRequestInfos(height, reqResponse)
			if err != nil {
				return errors.Wrap(err, "error while refreshing requests")
			}

			// Everything is ok
			return nil
		},
	}
}
