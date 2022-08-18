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

// oracleScriptsCmd returns a Cobra command that allows to refresh oracle scripts.
func oracleScriptsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "oracle-scripts",
		Short: "Refresh the information about oracle scripts taking them from the latest known height",
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

			// Get all oracle scripts
			var txs []*tmctypes.ResultTx

			// Firstly, MsgCreateOracleScript
			query := fmt.Sprintf("message.action='/oracle.v1.MsgCreateOracleScript'")
			createOracleScriptTxs, err := utils.QueryTxs(parseCtx.Node, query)
			if err != nil {
				return errors.Wrap(err, "Failed to get MsgCreateOracleScript messages")
			}

			txs = append(txs, createOracleScriptTxs...)

			// Then - MsgEditOracleScript
			query = fmt.Sprintf("message.action='/oracle.v1.MsgEditOracleScript'")
			editOracleScriptTxs, err := utils.QueryTxs(parseCtx.Node, query)
			if err != nil {
				return errors.Wrap(err, "Failed to get MsgEditOracleScript messages")
			}

			txs = append(txs, editOracleScriptTxs...)

			err = oracleModule.HandleOracleTxs(txs, parseCtx)
			if err != nil {
				return errors.Wrap(err, "Error while handling oracle module message")
			}

			// Everything is ok
			return nil
		},
	}
}
