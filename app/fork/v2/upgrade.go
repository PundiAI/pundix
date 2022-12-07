package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	pxtypes "github.com/pundix/pundix/types"
)

const (
	name = "pxv2"
	info = `'{"binaries":{"darwin/arm64":"https://github.com/pundix/pundix/releases/download/v0.2.0/pundix_0.2.0_Darwin_arm64.tar.gz","darwin/amd64":"https://github.com/pundix/pundix/releases/download/v0.2.0/pundix_0.2.0_Darwin_amd64.tar.gz","linux/arm64":"https://github.com/pundix/pundix/releases/download/v0.2.0/pundix_0.2.0_Linux_arm64.tar.gz","linux/amd64":"https://github.com/pundix/pundix/releases/download/v0.2.0/pundix_0.2.0_Linux_amd64.tar.gz","windows/x86_64":"https://github.com/pundix/pundix/releases/download/v0.2.0/pundix_0.2.0_Windows_x86_64.zip"}}'`
)

func Upgrade(ctx sdk.Context, upgradeKeeper upgradekeeper.Keeper) error {
	plan := types.Plan{
		Name:   name,
		Height: int64(pxtypes.V2SoftwareUpgradeHeight()),
		Info:   info,
	}
	ctx.Logger().With("fork/v2").Info("schedule upgrade begin", "plan", plan)
	err := upgradeKeeper.ScheduleUpgrade(ctx, plan)
	ctx.Logger().With("fork/v2").Info("schedule upgrade done", "err", err)
	return err
}
