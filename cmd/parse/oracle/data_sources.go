package oracle

import (
	"fmt"
	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/oracle"
	modulestypes "github.com/forbole/bdjuno/v3/modules/types"
	"github.com/forbole/bdjuno/v3/utils"
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/types/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// dataSourcesCmd returns a Cobra command that allows to refresh data sources.
func dataSourcesCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "data-sources",
		Short: "Refresh the information about data sources taking them from the latest known height",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return errors.Wrap(err, "Failed to get parser context")
			}

			sources, err := modulestypes.BuildSources(config.Cfg.Node, parseCtx.EncodingConfig)
			if err != nil {
				return errors.Wrap(err, "Failed to build sources")
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build the oracle module
			oracleModule := oracle.NewModule(sources.OracleSource, db)

			// Get all data sources
			var txs []*tmctypes.ResultTx

			// Firstly, MsgCreateDataSource
			query := fmt.Sprintf("message.action='/oracle.v1.MsgCreateDataSource'")
			createDataSourceTxs, err := utils.QueryTxs(parseCtx.Node, query)
			if err != nil {
				return errors.Wrap(err, "Failed to get MsgCreateDataSource messages")
			}

			txs = append(txs, createDataSourceTxs...)

			// Then - MsgEditDataSource
			query = fmt.Sprintf("message.action='/oracle.v1.MsgEditDataSource'")
			editDataSourceTxs, err := utils.QueryTxs(parseCtx.Node, query)
			if err != nil {
				return errors.Wrap(err, "Failed to get MsgEditDataSource messages")
			}

			txs = append(txs, editDataSourceTxs...)

			err = oracleModule.HandleOracleTxs(txs, parseCtx)
			if err != nil {
				return errors.Wrap(err, "Error while handling oracle module message")
			}

			// Everything is ok
			return nil
		},
	}
}
