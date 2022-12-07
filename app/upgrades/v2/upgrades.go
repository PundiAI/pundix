package v2

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/pundix/pundix/app/keepers"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v9
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, _ module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("start to run module migrations...")
		fromVM := map[string]uint64{
			"auth":         2,
			"bank":         2,
			"capability":   1,
			"crisis":       1,
			"distribution": 2,
			"evidence":     1,
			"genutil":      1,
			"gov":          2,
			"ibc":          2,
			"mint":         1,
			"other":        1,
			"params":       1,
			"slashing":     2,
			"staking":      2,
			"transfer":     1,
			"upgrade":      1,
			"vesting":      1,
		}
		// Leave modules are as-is to avoid running InitGenesis.
		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return vm, err
		}
		ctx.Logger().Info("running the rest of the upgrade handler...")

		// 1. change gov min deposit
		ChangeGovDepositParams(ctx, keepers.GovKeeper, keepers.StakingKeeper)

		// 2. change gov voting period
		ChangeGovVotingParams(ctx, keepers.GovKeeper)

		ctx.Logger().Info("upgrade complete")
		return vm, err
	}
}

func ChangeGovDepositParams(ctx sdk.Context, govKeeper govkeeper.Keeper, stakingKeeper stakingkeeper.Keeper) {
	bondDenom := stakingKeeper.BondDenom(ctx)

	govDepositParams := govKeeper.GetDepositParams(ctx)
	coinOne := sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	for i, coin := range govDepositParams.MinDeposit {
		if coin.Denom == bondDenom {
			govDepositParams.MinDeposit[i].Amount = coinOne.Mul(sdk.NewInt(GovMinDeposit))
			break
		}
	}
	ctx.Logger().Info("change x/gov module deposit params:minDeposit", "params", govDepositParams)
	govKeeper.SetDepositParams(ctx, govDepositParams)
}

func ChangeGovVotingParams(ctx sdk.Context, govKeeper govkeeper.Keeper) {
	govVotingParams := govKeeper.GetVotingParams(ctx)
	govVotingParams.VotingPeriod = GovVotingPeriod
	ctx.Logger().Info("change x/gov module voting params:votingPeriod", "params", govVotingParams)
	govKeeper.SetVotingParams(ctx, govVotingParams)
}
