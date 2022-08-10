package oracle

import (
	"errors"
	"fmt"
	modulestypes "github.com/forbole/bdjuno/v3/modules/types"

	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/oracle"
)

// requestsCmd returns a Cobra command that allows to refresh requests.
func requestsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "requests",
		Short: "Refresh the information about requests taking them from the latest known height",
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

			// Get all requests
			requests, err := sources.OracleSource.GetRequests(height)
			if err != nil {
				return fmt.Errorf("error while getting requests: %s", err)
			}

			//Refresh each request
			for _, request := range requests {
				err = oracleModule.RefreshRequestInfos(int64(request.Result.RequestID), height)
				if err != nil {
					return fmt.Errorf("error while refreshing requests: %s", err)
				}
			}

			//s := fmt.Sprintf("Norm?: %d", int(requests.OracleScriptID))
			return errors.New("Norm?: ")
		},
	}
}
