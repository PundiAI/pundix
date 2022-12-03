package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
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

		// do something

		ctx.Logger().Info("upgrade complete")
		return vm, err
	}
}
