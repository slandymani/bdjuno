package types

import (
	"fmt"
	"os"

	"cosmossdk.io/log"
	onfttypes "github.com/ODIN-PROTOCOL/odin-core/x/onft/types"
	wasmkeeper "github.com/ODIN-PROTOCOL/wasmd/x/wasm/keeper"
	wasmtypes "github.com/ODIN-PROTOCOL/wasmd/x/wasm/types"
	onftsource "github.com/forbole/callisto/v4/modules/onft/source"
	wasmsource "github.com/forbole/callisto/v4/modules/wasm/source"
	"github.com/forbole/juno/v6/node/remote"

	mintkeeper "github.com/ODIN-PROTOCOL/odin-core/x/mint/keeper"
	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	oraclekeeper "github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v6/node/local"

	nodeconfig "github.com/forbole/juno/v6/node/config"

	odinapp "github.com/ODIN-PROTOCOL/odin-core/app"
	banksource "github.com/forbole/callisto/v4/modules/bank/source"
	localbanksource "github.com/forbole/callisto/v4/modules/bank/source/local"
	remotebanksource "github.com/forbole/callisto/v4/modules/bank/source/remote"
	distrsource "github.com/forbole/callisto/v4/modules/distribution/source"
	localdistrsource "github.com/forbole/callisto/v4/modules/distribution/source/local"
	remotedistrsource "github.com/forbole/callisto/v4/modules/distribution/source/remote"
	govsource "github.com/forbole/callisto/v4/modules/gov/source"
	localgovsource "github.com/forbole/callisto/v4/modules/gov/source/local"
	remotegovsource "github.com/forbole/callisto/v4/modules/gov/source/remote"
	mintsource "github.com/forbole/callisto/v4/modules/mint/source"
	localmintsource "github.com/forbole/callisto/v4/modules/mint/source/local"
	remotemintsource "github.com/forbole/callisto/v4/modules/mint/source/remote"
	localonftsource "github.com/forbole/callisto/v4/modules/onft/source/local"
	remoteonftsource "github.com/forbole/callisto/v4/modules/onft/source/remote"
	oraclesource "github.com/forbole/callisto/v4/modules/oracle/source"
	localoraclesource "github.com/forbole/callisto/v4/modules/oracle/source/local"
	remoteoraclesource "github.com/forbole/callisto/v4/modules/oracle/source/remote"
	slashingsource "github.com/forbole/callisto/v4/modules/slashing/source"
	localslashingsource "github.com/forbole/callisto/v4/modules/slashing/source/local"
	remoteslashingsource "github.com/forbole/callisto/v4/modules/slashing/source/remote"
	stakingsource "github.com/forbole/callisto/v4/modules/staking/source"
	localstakingsource "github.com/forbole/callisto/v4/modules/staking/source/local"
	remotestakingsource "github.com/forbole/callisto/v4/modules/staking/source/remote"
	localwasmsource "github.com/forbole/callisto/v4/modules/wasm/source/local"
	remotewasmsource "github.com/forbole/callisto/v4/modules/wasm/source/remote"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
)

type Sources struct {
	BankSource     banksource.Source
	DistrSource    distrsource.Source
	GovSource      govsource.Source
	MintSource     mintsource.Source
	OracleSource   oraclesource.Source
	SlashingSource slashingsource.Source
	StakingSource  stakingsource.Source
	WasmSource     wasmsource.Source
	ONftSource     onftsource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, cdc codec.Codec) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return buildLocalSources(cfg, cdc)

	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

func buildLocalSources(cfg *local.Details, cdc codec.Codec) (*Sources, error) {
	source, err := local.NewSource(cfg.Home, cdc)
	if err != nil {
		return nil, err
	}

	appOpts := make(simtestutil.AppOptionsMap, 0)

	app := odinapp.NewOdinApp(
		log.NewLogger(os.Stdout), source.StoreDB, nil, false,
		map[int64]bool{}, appOpts, 0, nil,
	)

	sources := &Sources{
		BankSource:     localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
		DistrSource:    localdistrsource.NewSource(source, distrkeeper.NewQuerier(app.DistrKeeper)),
		GovSource:      localgovsource.NewSource(source, govkeeper.NewQueryServer(&app.GovKeeper)),
		MintSource:     localmintsource.NewSource(source, mintkeeper.NewQueryServerImpl(app.MintKeeper)),
		OracleSource:   localoraclesource.NewSource(source, oraclekeeper.Querier{Keeper: app.OracleKeeper}),
		SlashingSource: localslashingsource.NewSource(source, slashingtypes.QueryServer(app.SlashingKeeper)),
		StakingSource:  localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: app.StakingKeeper}),
		WasmSource:     localwasmsource.NewSource(source, wasmkeeper.Querier(&app.WasmKeeper)),
		ONftSource:     localonftsource.NewSource(source, onfttypes.QueryServer(app.ONFTKeeper)),
	}

	// Mount and initialize the stores
	err = source.MountKVStores(app, "keys")
	if err != nil {
		return nil, err
	}

	err = source.InitStores()
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func buildRemoteSources(cfg *remote.Details) (*Sources, error) {
	source, err := remote.NewSource(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote source: %s", err)
	}

	return &Sources{
		BankSource:     remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:    remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:      remotegovsource.NewSource(source, govtypesv1.NewQueryClient(source.GrpcConn)),
		MintSource:     remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		OracleSource:   remoteoraclesource.NewSource(source, oracletypes.NewQueryClient(source.GrpcConn)),
		SlashingSource: remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:  remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
		WasmSource:     remotewasmsource.NewSource(source, wasmtypes.NewQueryClient(source.GrpcConn)),
		ONftSource:     remoteonftsource.NewSource(source, onfttypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
