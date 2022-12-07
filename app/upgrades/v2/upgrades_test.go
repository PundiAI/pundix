package v2_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/pundix/pundix/app"
	v2 "github.com/pundix/pundix/app/upgrades/v2"
)

type UpgradeTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.PundixApp
}

func (suite *UpgradeTestSuite) SetupTest() {
	suite.app = app.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
		Height:  1,
	})
}

func TestUpgradeTestSuite(t *testing.T) {
	s := new(UpgradeTestSuite)
	suite.Run(t, s)
}

func (suite *UpgradeTestSuite) TestChangeGovDepositParams() {
	govDepositParams := suite.app.GovKeeper.GetDepositParams(suite.ctx)
	suite.Require().NotNil(govDepositParams)
	suite.Require().EqualValues(1, govDepositParams.MinDeposit.Len())
	bondDenom := suite.app.StakingKeeper.BondDenom(suite.ctx)
	suite.Require().NotEmpty(bondDenom)

	suite.Require().EqualValues(bondDenom, govDepositParams.MinDeposit[0].Denom)

	suite.Require().True(govDepositParams.MinDeposit[0].Amount.
		Equal(
			sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)).Mul(sdk.NewInt(10000)),
		))
	v2.ChangeGovDepositParams(suite.ctx, suite.app.GovKeeper, suite.app.StakingKeeper)

	newGovDepositParams := suite.app.GovKeeper.GetDepositParams(suite.ctx)
	suite.Require().NotNil(newGovDepositParams)
	suite.Require().EqualValues(1, newGovDepositParams.MinDeposit.Len())
	suite.Require().EqualValues(bondDenom, newGovDepositParams.MinDeposit[0].Denom)

	suite.Require().True(newGovDepositParams.MinDeposit[0].Amount.
		Equal(
			sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)).Mul(sdk.NewInt(v2.GovMinDeposit)),
		))
}

func (suite *UpgradeTestSuite) TestChangeGovVotingParams() {
	govVotingParams := suite.app.GovKeeper.GetVotingParams(suite.ctx)
	suite.Require().NotNil(govVotingParams)
	suite.Require().EqualValues(time.Hour*24*14, govVotingParams.VotingPeriod)
	v2.ChangeGovVotingParams(suite.ctx, suite.app.GovKeeper)

	newGovVotingParams := suite.app.GovKeeper.GetVotingParams(suite.ctx)
	suite.Require().NotNil(newGovVotingParams)
	suite.Require().EqualValues(v2.GovVotingPeriod, newGovVotingParams.VotingPeriod)
}
