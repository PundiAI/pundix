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
	coin, err := sdk.ParseCoinNormalized("1PUNDIX")
	assert.NoError(t, err)
	assert.Equal(t, sdk.NewCoin(StakingBondDenom(), sdk.NewInt(1e18)).String(), coin.String())
	t.Log(coin.String())
}

func TestParseGasPrices(t *testing.T) {
	SetChainId(MainnetChainId)
	coin, err := sdk.ParseCoinNormalized("0.000002PUNDIX")
	assert.NoError(t, err)
	assert.Equal(t, sdk.NewCoin(StakingBondDenom(), sdk.NewInt(2e12)).String(), coin.String())
	t.Log(coin.String())
}
