package app

import (
	"io"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	kudosmodule "github.com/pavlenkotm/cosmos-kudos-module/x/kudos"
	kudoskeeper "github.com/pavlenkotm/cosmos-kudos-module/x/kudos/keeper"
	kudostypes "github.com/pavlenkotm/cosmos-kudos-module/x/kudos/types"
)

// ExampleApp extends an ABCI application with the kudos module integrated
type ExampleApp struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	// keepers
	AccountKeeper authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	KudosKeeper   kudoskeeper.Keeper

	// module manager
	mm *module.Manager
}

// NewExampleApp returns a reference to an initialized ExampleApp
func NewExampleApp(
	logger log.Logger,
	db storetypes.CommitMultiStore,
	traceStore io.Writer,
	loadLatest bool,
	appOpts interface{},
	baseAppOptions ...func(*baseapp.BaseApp),
) *ExampleApp {
	interfaceRegistry, _ := types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
		ProtoFiles: nil,
		SigningOptions: nil,
	})

	appCodec := codec.NewProtoCodec(interfaceRegistry)
	legacyAmino := codec.NewLegacyAmino()

	bApp := baseapp.NewBaseApp("kudos-app", logger, db, nil, baseAppOptions...)

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		kudostypes.StoreKey,
	)

	app := &ExampleApp{
		BaseApp:           bApp,
		cdc:               legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
	}

	// Initialize keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		nil,
		sdk.Bech32MainPrefix,
		authtypes.NewModuleAddress("gov").String(),
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.AccountKeeper,
		nil,
		authtypes.NewModuleAddress("gov").String(),
		logger,
	)

	// Initialize Kudos Keeper
	app.KudosKeeper = kudoskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[kudostypes.StoreKey]),
		logger,
	)

	// Register modules
	app.mm = module.NewManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, nil),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, nil),
		kudosmodule.NewAppModule(appCodec, app.KudosKeeper),
	)

	// Set order of begin blockers, end blockers, and init genesis
	app.mm.SetOrderBeginBlockers(
		authtypes.ModuleName,
		banktypes.ModuleName,
		kudostypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		authtypes.ModuleName,
		banktypes.ModuleName,
		kudostypes.ModuleName,
	)

	app.mm.SetOrderInitGenesis(
		authtypes.ModuleName,
		banktypes.ModuleName,
		kudostypes.ModuleName,
	)

	// Register module routes and query routes
	app.mm.RegisterServices(module.NewConfigurator(appCodec, bApp.MsgServiceRouter(), bApp.GRPCQueryRouter()))

	// Mount stores
	if err := app.LoadLatestVersion(); err != nil {
		panic(err)
	}

	return app
}

// LoadLatestVersion loads the latest application version
func (app *ExampleApp) LoadLatestVersion() error {
	return app.LoadVersion(app.LastBlockHeight())
}

// LoadVersion loads the app at a specific version
func (app *ExampleApp) LoadVersion(version int64) error {
	return app.BaseApp.LoadVersion(version)
}

// LegacyAmino returns ExampleApp's amino codec
func (app *ExampleApp) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns ExampleApp's app codec
func (app *ExampleApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns ExampleApp's InterfaceRegistry
func (app *ExampleApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// ModuleManager returns the app's module manager
func (app *ExampleApp) ModuleManager() *module.Manager {
	return app.mm
}

// RegisterAPIRoutes registers all application module routes with the provided API server
func (app *ExampleApp) RegisterAPIRoutes(apiSvr interface{}, apiConfig interface{}) {
	// Register gRPC Gateway routes here
}

// InitChainer application update at chain initialization
func (app *ExampleApp) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState map[string]json.RawMessage
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		return nil, err
	}

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// BeginBlocker application updates every begin block
func (app *ExampleApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.mm.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *ExampleApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}
