package handlers

import (
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/forbole/callisto/v4/database"
	"github.com/forbole/callisto/v4/modules/actions/types"
	types2 "github.com/forbole/callisto/v4/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func GetVotingResult(ctx *types.Context, payload *types.Payload, db *database.Db) (interface{}, error) {
	log.Debug().Msg("executing voting result action")

	proposalID := payload.GetID()

	votingResultCountResponse, err := getVotingCountResult(proposalID, db)
	if err != nil {
		return nil, err
	}

	votingResultWeightedResponse, err := getVotingWeightedResult(proposalID, db)
	if err != nil {
		return nil, err
	}

	return VotingResultResponse{
		CountResponse:    votingResultCountResponse,
		WeightedResponse: votingResultWeightedResponse,
	}, nil
}

func getVotingWeightedResult(id int64, db *database.Db) (VotingResultWeightedResponse, error) {
	var response []types2.TallyResult
	stmt := `SELECT * FROM proposal_tally_result WHERE proposal_id = $1`

	err := db.Sqlx.Select(&response, stmt, id)
	if err != nil {
		return VotingResultWeightedResponse{}, errors.Wrap(err, "failed to select proposal tally result")
	}

	if len(response) == 0 {
		return VotingResultWeightedResponse{}, nil
	}

	return VotingResultWeightedResponse{
		Yes:        response[0].Yes,
		Abstain:    response[0].Abstain,
		No:         response[0].No,
		NoWithVeto: response[0].NoWithVeto,
	}, nil
}

func getVotingCountResult(id int64, db *database.Db) (VotingResultCountResponse, error) {
	var options []string
	var response VotingResultCountResponse

	stmt := `SELECT option FROM proposal_vote WHERE proposal_id = $1`
	err := db.Sqlx.Select(&options, stmt, id)
	if err != nil {
		return VotingResultCountResponse{}, errors.Wrap(err, "failed to select proposal votes")
	}

	for _, option := range options {
		switch govtypesv1.VoteOption_value[option] {
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
			return VotingResultCountResponse{}, errors.New("wrong voting option type occurred")
		}
	}
	response.TotalCount = int64(len(options))

	return response, nil
}

type VotingResultResponse struct {
	CountResponse    VotingResultCountResponse    `json:"count"`
	WeightedResponse VotingResultWeightedResponse `json:"weighted"`
}

type VotingResultWeightedResponse struct {
	Yes        string `json:"yes"`
	No         string `json:"no"`
	NoWithVeto string `json:"no_with_veto"`
	Abstain    string `json:"abstain"`
}

type VotingResultCountResponse struct {
	Yes        int64 `json:"yes"`
	No         int64 `json:"no"`
	NoWithVeto int64 `json:"no_with_veto"`
	Abstain    int64 `json:"abstain"`
	TotalCount int64 `json:"total_count"`
}
