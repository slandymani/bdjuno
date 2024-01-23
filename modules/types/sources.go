package types

import (
	"fmt"
	"os"

	"cosmossdk.io/simapp"
	"cosmossdk.io/simapp/params"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/forbole/juno/v5/node/remote"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	oraclekeeper "github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	telemetrytypes "github.com/ODIN-PROTOCOL/odin-core/x/telemetry/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v5/node/local"

	nodeconfig "github.com/forbole/juno/v5/node/config"

	odinapp "github.com/ODIN-PROTOCOL/odin-core/app"
	banksource "github.com/forbole/bdjuno/v4/modules/bank/source"
	localbanksource "github.com/forbole/bdjuno/v4/modules/bank/source/local"
	remotebanksource "github.com/forbole/bdjuno/v4/modules/bank/source/remote"
	distrsource "github.com/forbole/bdjuno/v4/modules/distribution/source"
	localdistrsource "github.com/forbole/bdjuno/v4/modules/distribution/source/local"
	remotedistrsource "github.com/forbole/bdjuno/v4/modules/distribution/source/remote"
	govsource "github.com/forbole/bdjuno/v4/modules/gov/source"
	localgovsource "github.com/forbole/bdjuno/v4/modules/gov/source/local"
	remotegovsource "github.com/forbole/bdjuno/v4/modules/gov/source/remote"
	mintsource "github.com/forbole/bdjuno/v4/modules/mint/source"
	localmintsource "github.com/forbole/bdjuno/v4/modules/mint/source/local"
	remotemintsource "github.com/forbole/bdjuno/v4/modules/mint/source/remote"
	oraclesource "github.com/forbole/bdjuno/v4/modules/oracle/source"
	localoraclesource "github.com/forbole/bdjuno/v4/modules/oracle/source/local"
	remoteoraclesource "github.com/forbole/bdjuno/v4/modules/oracle/source/remote"
	slashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source"
	localslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/local"
	remoteslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/remote"
	stakingsource "github.com/forbole/bdjuno/v4/modules/staking/source"
	localstakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/local"
	remotestakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/remote"
	telemetrysource "github.com/forbole/bdjuno/v4/modules/telemetry/source"
	localtelemetrysource "github.com/forbole/bdjuno/v4/modules/telemetry/source/local"
	remotetelemetrysource "github.com/forbole/bdjuno/v4/modules/telemetry/source/remote"
)

type Sources struct {
	BankSource      banksource.Source
	DistrSource     distrsource.Source
	GovSource       govsource.Source
	MintSource      mintsource.Source
	OracleSource    oraclesource.Source
	SlashingSource  slashingsource.Source
	StakingSource   stakingsource.Source
	TelemetrySource telemetrysource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, encodingConfig *params.EncodingConfig) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return buildLocalSources(cfg, encodingConfig)

	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

func buildLocalSources(cfg *local.Details, encodingConfig *params.EncodingConfig) (*Sources, error) {
	source, err := local.NewSource(cfg.Home, encodingConfig)
	if err != nil {
		return nil, err
	}

	app := odinapp.NewOdinApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, false, map[int64]bool{}, cfg.Home, 0,
		odinapp.MakeEncodingConfig(), simapp.EmptyAppOptions{}, false, 0,
	)

	sources := &Sources{
		BankSource:     localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
		DistrSource:    localdistrsource.NewSource(source, distrtypes.QueryServer(app.DistrKeeper)),
		GovSource:      localgovsource.NewSource(source, govtypesv1.QueryServer(app.GovKeeper)),
		MintSource:     localmintsource.NewSource(source, minttypes.QueryServer(app.MintKeeper)),
		OracleSource:    localoraclesource.NewSource(source, oraclekeeper.Querier{Keeper: app.OracleKeeper}),
		SlashingSource: localslashingsource.NewSource(source, slashingtypes.QueryServer(app.SlashingKeeper)),
		StakingSource:  localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: app.StakingKeeper}),
		TelemetrySource: localtelemetrysource.NewSource(source, telemetrytypes.QueryServer(app.TelemetryKeeper)),
	}

	// Mount and initialize the stores
	err = source.MountKVStores(app, "keys")
	if err != nil {
		return nil, err
	}

	err = source.MountTransientStores(app, "tkeys")
	if err != nil {
		return nil, err
	}

	err = source.MountMemoryStores(app, "memKeys")
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
		BankSource:      remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:     remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:       remotegovsource.NewSource(source, govtypesv1.NewQueryClient(source.GrpcConn)),
		MintSource:      remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		OracleSource:    remoteoraclesource.NewSource(source, oracletypes.NewQueryClient(source.GrpcConn)),
		SlashingSource:  remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:   remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
		TelemetrySource: remotetelemetrysource.NewSource(source, telemetrytypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
