package handlers

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/actions/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func GetVotingResult(ctx *types.Context, payload *types.Payload, db *database.Db) (interface{}, error) {
	log.Debug().Msg("executing voting result action")

	proposalID := payload.GetID()
	var options []string
	var response VotingResultResponse

	stmt := `SELECT option FROM proposal_vote WHERE proposal_id = $1`
	err := db.Sqlx.Select(&options, stmt, proposalID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select proposal votes")
	}

	for _, option := range options {
		switch govtypes.VoteOption_value[option] {
		case 0:
			continue
		case 1:
			response.Yes++
		case 2:
			response.Abstain++
		case 3:
			response.No++
		case 4:
			response.NoWithVeto++
		default:
			return nil, errors.New("wrong voting option type occured")
		}
	}
	response.TotalCount = int64(len(options))

	return response, nil
}

type VotingResultResponse struct {
	Yes        int64 `json:"yes"`
	No         int64 `json:"no"`
	NoWithVeto int64 `json:"no_with_veto"`
	Abstain    int64 `json:"abstain"`
	TotalCount int64 `json:"total_count"`
}
