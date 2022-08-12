package oracle

import (
	"encoding/hex"
	"fmt"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
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

func alternativeRequestCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "alternative-request [id]",
		Args:  cobra.ExactValidArgs(1),
		Short: "Refresh the information about selected request taking it from the latest known height",
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

			// Get specific request by id:
			var tx *tmctypes.ResultTx

			query := fmt.Sprintf("request.id='%s'", args[0])
			requestTx, err := utils.QueryTxs(parseCtx.Node, query)
			if err != nil {
				return err
			}

			if len(requestTx) != 1 {
				return errors.New(fmt.Sprintf("Request with id '%s' does not exist", args[0]))
			}

			tx = requestTx[0]

			transaction, err := parseCtx.Node.Tx(hex.EncodeToString(tx.Tx.Hash()))
			if err != nil {
				return err
			}

			// Handle only MsgRequestData instance
			for index, msg := range transaction.GetMsgs() {
				_, isMsgReqData := msg.(*oracletypes.MsgRequestData)
				if !isMsgReqData {
					continue
				}

				err = oracleModule.HandleMsg(index, msg, transaction)
				if err != nil {
					return fmt.Errorf("error while handling oracle module message: %s", err)
				}
			}

			// Everything is ok
			return nil
		},
	}
}
