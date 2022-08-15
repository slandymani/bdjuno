package oracle

import (
	"fmt"
	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/oracle"
	modulestypes "github.com/forbole/bdjuno/v3/modules/types"
	"github.com/forbole/bdjuno/v3/utils"
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/types/config"
	"github.com/spf13/cobra"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

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

			// Get all requests
			var txs []*tmctypes.ResultTx

			query := fmt.Sprintf("message.action='/oracle.v1.MsgRequestData'")
			requestsTx, err := utils.QueryTxs(parseCtx.Node, query)
			if err != nil {
				return err
			}

			txs = append(txs, requestsTx...)

			//TODO: REQ REPORTS

			err = oracleModule.HandleOracleTxs(txs, parseCtx)
			if err != nil {
				return fmt.Errorf("error while handling oracle module message: %s", err)
			}

			// Everything is ok
			return nil
		},
	}
}
