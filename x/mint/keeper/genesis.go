package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// InitGenesis new mint genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, ak minttypes.AccountKeeper, data *minttypes.GenesisState) {
	keeper.Keeper.SetMinter(ctx, data.Minter)
	keeper.SetParams(ctx, data.Params)
	ak.GetModuleAccount(ctx, minttypes.ModuleName)
}
