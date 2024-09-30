package handlers

import (
	"reflect"
	"strings"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	mint1 "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc1 "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/forbole/callisto/v4/database"
	"github.com/forbole/callisto/v4/modules/actions/types"
)

type Params struct {
	Key       string `json:"key"`
	ValueType string `json:"value_type"`
}

type ModuleParamsInfo struct {
	ModuleName string
	ParamPairs paramtypes.ParamSetPairs
}

type VoteParams struct {
	Module string   `json:"module"`
	Params []Params `json:"params"`
}

func GetVoteProposals(_ *types.Context, _ *types.Payload, _ *database.Db) (interface{}, error) {
	var typesMap = map[string]string{
		"types.Coins":           "[]{\"amount\": string, \"denom\": string}",
		"[]types.Exchange":      "[]{\"from\": string, \"to\": string, \"rate_multiplier\": types.Dec}",
		"[]*types.SendEnabled":  "[]{\"denom\": string, \"enabled\": bool}",
		"types.RewardThreshold": "\"amount\": []{\"amount\": string, \"denom\": string}, \"blocks\": uint64}",
		"[]*types.AllowedDenom": "[]{\"token_unit_denom\": string, \"token_denom\": string}",
	}

	keyParams := make(map[string][]Params)

	oracleParams := oracletypes.DefaultParams()
	mintParams := minttypes.DefaultParams()
	authParams := authtypes.DefaultParams()
	bankParams := banktypes.DefaultParams()
	distributionParams := distributiontypes.DefaultParams()
	ibcParams := ibctypes.DefaultParams()
	ibc1Params := ibc1.DefaultParams()
	mint1Params := mint1.DefaultParams()
	slashingParams := slashingtypes.DefaultParams()
	stakingParams := stakingtypes.DefaultParams()

	modulesParamsInfos := []ModuleParamsInfo{
		{
			ModuleName: "oracle",
			ParamPairs: oracleParams.ParamSetPairs(),
		},
		{
			ModuleName: "mint",
			ParamPairs: mintParams.ParamSetPairs(),
		},
		{
			ModuleName: "auth",
			ParamPairs: authParams.ParamSetPairs(),
		},
		{
			ModuleName: "bank",
			ParamPairs: bankParams.ParamSetPairs(),
		},
		{
			ModuleName: "distribution",
			ParamPairs: distributionParams.ParamSetPairs(),
		},
		{
			ModuleName: "ibc",
			ParamPairs: ibc1Params.ParamSetPairs(),
		},
		{
			ModuleName: "ibc",
			ParamPairs: ibcParams.ParamSetPairs(),
		},
		{
			ModuleName: "mint",
			ParamPairs: mint1Params.ParamSetPairs(),
		},
		{
			ModuleName: "slashing",
			ParamPairs: slashingParams.ParamSetPairs(),
		},
		{
			ModuleName: "staking",
			ParamPairs: stakingParams.ParamSetPairs(),
		},
	}

	for _, param := range modulesParamsInfos {
		moduleParams := GetModuleParams(param.ParamPairs, typesMap)

		if keyParams[param.ModuleName] != nil {
			keyParams[param.ModuleName] = append(keyParams[param.ModuleName], moduleParams...)
		} else {
			keyParams[param.ModuleName] = moduleParams
		}
	}

	keyParams["gov"] = []Params{
		{
			Key:       "depositparams",
			ValueType: "{\"min_deposit\": []{\"amount\": string, \"denom\": string}, \"max_deposit_period\": time.Duration}",
		},
		{
			Key:       "votingparams",
			ValueType: "{\"voting_period\": time.Duration}",
		},
		{
			Key:       "tallyparams",
			ValueType: "{\"quorum\": types.Dec, \"threshold\": types.Dec, \"veto_threshold\": types.Dec}",
		},
	}

	keyParams["crisis"] = []Params{
		{
			Key:       "ConstantFee",
			ValueType: "{\"amount\": string, \"denom\": string}",
		},
	}

	response := make([]VoteParams, 0)

	for module, params := range keyParams {
		response = append(response, VoteParams{
			Module: module,
			Params: params,
		})
	}

	return response, nil
}

func GetModuleParams(moduleParamsPairs paramtypes.ParamSetPairs, typesMap map[string]string) []Params {
	params := make([]Params, 0)

	for i := 0; i < len(moduleParamsPairs); i++ {
		valueType := reflect.TypeOf(moduleParamsPairs[i].Value).String()
		valueType = strings.Trim(valueType, "*")

		if typesMap[valueType] != "" {
			valueType = typesMap[valueType]
		}

		params = append(params, Params{
			Key:       string(moduleParamsPairs[i].Key),
			ValueType: valueType,
		})
	}

	return params
}
