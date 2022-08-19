package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestSetChainId(t *testing.T) {
	SetChainId(MainnetChainId)
	assert.Equal(t, ChainId(), MainnetChainId)
	SetChainId(TestnetChainId)
	assert.Equal(t, ChainId(), MainnetChainId)
}

func TestRegisterDenom(t *testing.T) {
	SetChainId(MainnetChainId)
	coin, err := sdk.ParseCoinsNormalized("1PUNDIX")
	assert.NoError(t, err)
	assert.Equal(t, sdk.NewCoins(sdk.NewCoin(StakingBondDenom(), sdk.NewInt(1e18))), coin)
}
