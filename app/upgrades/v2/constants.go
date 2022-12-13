package v2

import (
	"time"

	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"

	"github.com/pundix/pundix/app/upgrades"
)

const (
	// UpgradeName is the shared upgrade plan name for mainnet
	UpgradeName = "pxv2"
)

const (
	GovVotingPeriod = time.Hour * 24 * 7
	GovMinDeposit   = 3000
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
