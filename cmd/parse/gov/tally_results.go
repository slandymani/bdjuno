package gov

import (
	"fmt"

	govsource "github.com/forbole/callisto/v4/modules/gov/source"
	modulestypes "github.com/forbole/callisto/v4/modules/types"
	"github.com/forbole/callisto/v4/types"
	"github.com/forbole/callisto/v4/utils"

	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	parsecmdtypes "github.com/forbole/juno/v6/cmd/parse/types"
	"github.com/forbole/juno/v6/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/callisto/v4/database"
)

// tallyResultsCmd returns the Cobra command allowing to fix all things related to a tally results
func tallyResultsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "tally-results",
		Short: "Updates tally results of all proposals",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			codec := utils.GetCodec()

			sources, err := modulestypes.BuildSources(config.Cfg.Node, codec)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			err = refreshTallyResults(db, sources.GovSource)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func refreshTallyResults(db *database.Db, source govsource.Source) error {
	block, err := db.GetLastBlock()
	if err != nil {
		return fmt.Errorf("error while getting last block: %s", err)
	}

	proposalIDs, err := db.GetAllProposalsIds()
	if err != nil {
		return fmt.Errorf("error while getting proposal ids: %s", err)
	}

	tallys := make([]types.TallyResult, len(proposalIDs))
	for i, id := range proposalIDs {
		tally, err := source.TallyResult(block.Height, id)
		if err != nil {
			return fmt.Errorf("error while getting tally results: %s", err)
		}

		tallys[i] = parseTally(id, block.Height, tally)
	}

	return db.SaveTallyResults(tallys)
}

func parseTally(id uint64, height int64, tally *govtypesv1.TallyResult) types.TallyResult {
	return types.TallyResult{
		ProposalID: id,
		Yes:        tally.YesCount,
		Abstain:    tally.AbstainCount,
		No:         tally.NoCount,
		NoWithVeto: tally.NoWithVetoCount,
		Height:     height,
	}
}
