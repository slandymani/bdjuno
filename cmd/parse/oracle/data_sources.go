package oracle

import (
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/forbole/callisto/v4/database"
	"github.com/forbole/callisto/v4/modules/oracle"
	modulestypes "github.com/forbole/callisto/v4/modules/types"
	"github.com/forbole/callisto/v4/utils"
	parsecmdtypes "github.com/forbole/juno/v6/cmd/parse/types"
	"github.com/forbole/juno/v6/types/config"
	"github.com/pkg/errors"
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
				return errors.Wrap(err, "Failed to get parser context")
			}

			codec := utils.GetCodec()

			sources, err := modulestypes.BuildSources(config.Cfg.Node, codec)
			if err != nil {
				return errors.Wrap(err, "Failed to build sources")
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build the oracle module
			oracleModule := oracle.NewModule(sources.OracleSource, db, codec)

			// Get all data sources
			var txs []*tmctypes.ResultTx

			// Firstly, MsgCreateDataSource
			query := "message.action='/oracle.v1.MsgCreateDataSource'"
			createDataSourceTxs, err := utils.QueryTxs(parseCtx.Node, query)
			if err != nil {
				return errors.Wrap(err, "Failed to get MsgCreateDataSource messages")
			}

			txs = append(txs, createDataSourceTxs...)

			// Then - MsgEditDataSource
			query = "message.action='/oracle.v1.MsgEditDataSource'"
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
