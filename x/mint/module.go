package mint

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/pundix/pundix/x/mint/types"

	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/pundix/pundix/x/mint/keeper"
)

var _ module.AppModule = AppModule{}

// AppModule implements an application module for the mint module.
type AppModule struct {
	mint.AppModule
	mint.AppModuleBasic

	keeper     keeper.Keeper
	authKeeper minttypes.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, ak minttypes.AccountKeeper) AppModule {
	appModule := mint.NewAppModule(cdc, keeper.Keeper, ak)
	return AppModule{
		keeper:         keeper,
		authKeeper:     ak,
		AppModule:      appModule,
		AppModuleBasic: appModule.AppModuleBasic,
	}
}

// ValidateGenesis performs genesis state validation for the mint module.
func (AppModule) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data minttypes.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", minttypes.ModuleName, err)
	}

	return types.ValidateGenesis(data)
}

// Name returns the mint module's name.
func (AppModule) Name() string {
	return minttypes.ModuleName
}

// LegacyQuerierHandler returns the mint module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return mintkeeper.NewQuerier(am.keeper.Keeper, legacyQuerierCdc)
}

// InitGenesis performs genesis initialization for the mint module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState minttypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	keeper.InitGenesis(ctx, am.keeper, am.authKeeper, &genesisState)
	return []abci.ValidatorUpdate{}
}

// BeginBlock returns the begin blocker for the mint module.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	mint.BeginBlocker(ctx, am.keeper.Keeper)
}
