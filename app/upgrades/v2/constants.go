package v2

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"

	"github.com/pundix/pundix/app/upgrades"
)

const (
	// UpgradeName is the shared upgrade plan name for mainnet
	UpgradeName = "v0.2.0"
	// UpgradeInfo defines the binaries that will be used for the upgrade
	UpgradeInfo = `'{"binaries":{"darwin/arm64":"https://github.com/pundix/pundix/releases/download/v0.2.0/pundix_0.2.0_Darwin_arm64.tar.gz","darwin/amd64":"https://github.com/pundix/pundix/releases/download/v0.2.0/pundix_0.2.0_Darwin_amd64.tar.gz","linux/arm64":"https://github.com/pundix/pundix/releases/download/v0.2.0/pundix_0.2.0_Linux_arm64.tar.gz","linux/amd64":"https://github.com/pundix/pundix/releases/download/v0.2.0/pundix_0.2.0_Linux_amd64.tar.gz","windows/x86_64":"https://github.com/pundix/pundix/releases/download/v0.2.0/pundix_0.2.0_Windows_x86_64.zip"}}'`
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			authz.ModuleName,
			feegrant.ModuleName,
		},
	},
}
