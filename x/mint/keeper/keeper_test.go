package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/pundix/pundix/app"
	pxmintkeeper "github.com/pundix/pundix/x/mint/keeper"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.PundixApp
}

var s *KeeperTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(KeeperTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = app.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
		Height:  1,
	})
}

func (suite *KeeperTestSuite) TestSetParams() {
	tests := []struct {
		name    string
		params  minttypes.Params
		expPass bool
	}{
		{name: "pass - default params", params: minttypes.DefaultParams(), expPass: true},

		{name: "error - empty mint denom", params: minttypes.NewParams("", sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), 0), expPass: false},
		{name: "error - negative InflationRateChange", params: minttypes.NewParams("stake", sdk.NewDec(-1), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), 0), expPass: false},
		{name: "error - negative InflationMax", params: minttypes.NewParams("stake", sdk.ZeroDec(), sdk.NewDec(-1), sdk.ZeroDec(), sdk.ZeroDec(), 0), expPass: false},
		{name: "error - negative InflationMin", params: minttypes.NewParams("stake", sdk.ZeroDec(), sdk.ZeroDec(), sdk.NewDec(-1), sdk.ZeroDec(), 0), expPass: false},
		{name: "error - negative GoalBonded", params: minttypes.NewParams("stake", sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.NewDec(-1), 0), expPass: false},
		{name: "error - zero blocksPerYear", params: minttypes.NewParams("stake", sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), 0), expPass: false},

		{name: "error - GoalBonded more then one", params: minttypes.NewParams("stake", sdk.OneDec(), sdk.OneDec(), sdk.OneDec(), sdk.NewDecWithPrec(11, 1), 1), expPass: false},
		{name: "pass - InflationRateChange more then one", params: minttypes.NewParams("stake", sdk.OneDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.OneDec(), 1), expPass: true},
		{name: "pass - InflationRateChange InflationMax InflationMin grate one", params: minttypes.NewParams("stake", sdk.NewDec(40), sdk.NewDec(40), sdk.NewDec(20), sdk.NewDecWithPrec(51, 2), minttypes.DefaultParams().BlocksPerYear), expPass: true},
	}
	for _, tt := range tests {
		suite.Run(fmt.Sprintf("Case %s", tt.name), func() {
			genesisState := minttypes.NewGenesisState(minttypes.DefaultInitialMinter(), tt.params)
			fn := require.Panics
			if tt.expPass {
				fn = require.NotPanics
			}

			fn(suite.T(), func() {
				pxmintkeeper.InitGenesis(suite.ctx, suite.app.MintKeeper, suite.app.AccountKeeper, genesisState)
				suite.app.MintKeeper.SetParams(suite.ctx, tt.params)
			})

			if tt.expPass {
				actualParams := suite.app.MintKeeper.Keeper.GetParams(suite.ctx)
				equalMintParams(suite.T(), tt.params, actualParams)
			}
		})
	}
}

func equalMintParams(t *testing.T, expect minttypes.Params, actual minttypes.Params) {
	require.Equal(t, expect.MintDenom, actual.MintDenom)
	require.Equal(t, expect.BlocksPerYear, actual.BlocksPerYear)
	require.True(t, expect.InflationRateChange.Equal(actual.InflationRateChange))
	require.True(t, expect.InflationMax.Equal(actual.InflationMax))
	require.True(t, expect.InflationMin.Equal(actual.InflationMin))
	require.True(t, expect.GoalBonded.Equal(actual.GoalBonded))
}
