package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	"github.com/stretchr/testify/suite"

	"github.com/pundix/pundix/app"
	pxibctesting "github.com/pundix/pundix/ibc/testing"
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainPx     *ibctesting.TestChain
	chainFx     *ibctesting.TestChain
	chainCosmos *ibctesting.TestChain

	queryClient transfertypes.QueryClient
}

var s *KeeperTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(KeeperTestSuite)
	suite.Run(t, s)
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.DoSetupTest(suite.T())
}

func (suite *KeeperTestSuite) DoSetupTest(t *testing.T) {
	suite.coordinator = pxibctesting.NewCoordinator(t, 2, 1)
	suite.chainPx = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainFx = suite.coordinator.GetChain(ibctesting.GetChainID(2))
	suite.chainCosmos = suite.coordinator.GetChain(ibctesting.GetChainID(3))

	queryHelper := baseapp.NewQueryServerTestHelper(suite.chainPx.GetContext(), suite.GetApp(suite.chainPx.App).InterfaceRegistry())
	transfertypes.RegisterQueryServer(queryHelper, suite.GetApp(suite.chainPx.App).TransferKeeper)
	suite.queryClient = transfertypes.NewQueryClient(queryHelper)
}

func (suite *KeeperTestSuite) GetApp(testingApp ibctesting.TestingApp) *app.PundixApp {
	pundixApp, ok := testingApp.(*app.PundixApp)
	suite.Require().True(ok)
	return pundixApp
}
