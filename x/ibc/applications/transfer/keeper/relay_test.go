package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcgotesting "github.com/cosmos/ibc-go/v3/testing"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/rand"

	"github.com/pundix/pundix/app"
	pxtypes "github.com/pundix/pundix/types"
	pundixtransfer "github.com/pundix/pundix/x/ibc/applications/transfer"
	pxtransfertypes "github.com/pundix/pundix/x/ibc/applications/transfer/types"
)

func (suite *KeeperTestSuite) TestOnRecvPacket() {
	baseDenom := "stake"
	senderAddr := sdk.AccAddress([]byte{0x1})
	receiveAddr := sdk.AccAddress([]byte{0x2})
	ibcDenomTrace := transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: baseDenom,
	}
	mintDenomTrace := transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: pxtypes.MintDenom(),
	}
	transferAmount := sdk.NewInt(100)
	testCases := []struct {
		name         string
		malleate     func(pxIbcTransferMsg *channeltypes.Packet)
		expPass      bool
		errorStr     string
		checkBalance bool
		expCoins     sdk.Coins
	}{
		{
			"pass - normal - ibc transfer packet",
			func(packet *channeltypes.Packet) {
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), transferAmount)),
		},
		{
			"pass - normal - px ibc transfer packet",
			func(packet *channeltypes.Packet) {
				packetData := pxtransfertypes.FungibleTokenPacketData{}
				pxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Router = ""
				packetData.Fee = sdk.ZeroInt().String()
				packet.Data = packetData.GetBytes()
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), transferAmount)),
		},
		{
			"pass - normal - ibc mint token - router is empty",
			func(packet *channeltypes.Packet) {
				packetData := pxtransfertypes.FungibleTokenPacketData{}
				pxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Denom = mintDenomTrace.GetFullDenomPath()
				packetData.Router = ""
				packetData.Fee = sdk.ZeroInt().String()
				packet.Data = packetData.GetBytes()
				// mint coin to channel escrowAddress
				mintCoin := sdk.NewCoins(sdk.NewCoin(pxtypes.MintDenom(), transferAmount))
				suite.Require().NoError(suite.GetApp(suite.chainPx.App).BankKeeper.MintCoins(suite.chainPx.GetContext(), transfertypes.ModuleName, mintCoin))
				escrowAddress := transfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
				suite.Require().NoError(suite.GetApp(suite.chainPx.App).BankKeeper.SendCoinsFromModuleToAccount(suite.chainPx.GetContext(), transfertypes.ModuleName, escrowAddress, sdk.NewCoins(sdk.NewCoin(pxtypes.MintDenom(), transferAmount))))
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(pxtypes.MintDenom(), transferAmount)),
		},
		{
			"error - not support router",
			func(packet *channeltypes.Packet) {
				packetData := pxtransfertypes.FungibleTokenPacketData{}
				pxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Denom = ibcDenomTrace.IBCDenom()
				packetData.Router = rand.Str(4)
				packet.Data = packetData.GetBytes()
			},
			false,
			// not support router error code: 101
			"ABCI code: 101: error handling packet on destination chain: see events for details",
			true,
			sdk.NewCoins(),
		},
		{
			"error - not support fee",
			func(packet *channeltypes.Packet) {
				packetData := pxtransfertypes.FungibleTokenPacketData{}
				pxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Denom = ibcDenomTrace.IBCDenom()
				packetData.Router = rand.Str(4)
				packetData.Fee = sdk.NewInt(100).String()
				packet.Data = packetData.GetBytes()
			},
			false,
			// not support router error code: 101, if router is empty fee -> empty
			"ABCI code: 101: error handling packet on destination chain: see events for details",
			true,
			sdk.NewCoins(),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			app := suite.GetApp(suite.chainPx.App)
			transferIBCModule := transfer.NewIBCModule(app.TransferKeeper)
			pundixIBCMiddleware := pundixtransfer.NewIBCModule(app.PundixTransferKeeper, transferIBCModule)
			packetData := transfertypes.NewFungibleTokenPacketData(baseDenom, transferAmount.String(), senderAddr.String(), receiveAddr.String())
			// only use timeout height
			packet := channeltypes.NewPacket(packetData.GetBytes(), 1, ibcgotesting.TransferPort, "channel-0", ibcgotesting.TransferPort, "channel-0", clienttypes.Height{
				RevisionNumber: 100,
				RevisionHeight: 100000,
			}, 0)
			tc.malleate(&packet)

			ackI := pundixIBCMiddleware.OnRecvPacket(suite.chainPx.GetContext(), packet, nil)
			suite.Require().NotNil(ackI)

			ack, ok := ackI.(channeltypes.Acknowledgement)

			if tc.expPass {
				suite.Require().Truef(ack.Success(), "error:%s,packetData:%s", ack.GetError(), string(packet.GetData()))
			} else {
				suite.Require().False(ack.Success())
				suite.Require().True(ok)
				suite.Require().Equalf(tc.errorStr, ack.GetError(), "packetData:%s", string(packet.GetData()))
			}

			if tc.checkBalance {
				bankKeeper := suite.GetApp(suite.chainPx.App).BankKeeper
				receiveAddrCoins := bankKeeper.GetAllBalances(suite.chainPx.GetContext(), receiveAddr)
				suite.Require().True(tc.expCoins.IsEqual(receiveAddrCoins), "exp:%s,actual:%s", tc.expCoins, receiveAddrCoins)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOnAcknowledgementPacket() {
	baseDenom := "stake"
	senderAddr := sdk.AccAddress([]byte{0x1})
	receiveAddr := sdk.AccAddress([]byte{0x2})
	transferAmount := sdk.NewInt(100)
	testCases := []struct {
		name         string
		malleate     func(pxIbcTransferMsg *channeltypes.Packet, ack *channeltypes.Acknowledgement)
		expPass      bool
		errorStr     string
		checkBalance bool
		expCoins     sdk.Coins
	}{
		{
			"pass - success ack - ibc transfer packet",
			func(packet *channeltypes.Packet, ack *channeltypes.Acknowledgement) {
				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				mintCoin(suite.T(), suite.chainPx.GetContext(), suite.GetApp(suite.chainPx.App), escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)))
			},
			true,
			"",
			true,
			sdk.NewCoins(),
		},
		{
			"pass - error ack - ibc transfer packet",
			func(packet *channeltypes.Packet, ack *channeltypes.Acknowledgement) {
				*ack = channeltypes.NewErrorAcknowledgement("test")

				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				mintCoin(suite.T(), suite.chainPx.GetContext(), suite.GetApp(suite.chainPx.App), escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)))
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			chain := suite.GetApp(suite.chainPx.App)
			transferIBCModule := transfer.NewIBCModule(chain.TransferKeeper)
			pundixIBCMiddleware := pundixtransfer.NewIBCModule(chain.PundixTransferKeeper, transferIBCModule)
			packetData := transfertypes.NewFungibleTokenPacketData(baseDenom, transferAmount.String(), senderAddr.String(), receiveAddr.String())
			// only use timeout height
			packet := channeltypes.NewPacket(packetData.GetBytes(), 1, ibcgotesting.TransferPort, "channel-0", ibcgotesting.TransferPort, "channel-0", clienttypes.Height{
				RevisionNumber: 100,
				RevisionHeight: 100000,
			}, 0)

			ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
			tc.malleate(&packet, &ack)

			err := pundixIBCMiddleware.OnAcknowledgementPacket(suite.chainPx.GetContext(), packet, ack.Acknowledgement(), nil)
			if tc.expPass {
				suite.Require().NoError(err, "packetData:%s", string(packet.GetData()))
			} else {
				suite.Require().NotNil(err)
				suite.Require().Equalf(tc.errorStr, err.Error(), "packetData:%s", string(packet.GetData()))
			}

			if tc.checkBalance {
				bankKeeper := suite.GetApp(suite.chainPx.App).BankKeeper
				senderAddrCoins := bankKeeper.GetAllBalances(suite.chainPx.GetContext(), senderAddr)
				suite.Require().True(tc.expCoins.IsEqual(senderAddrCoins), "exp:%s,actual:%s", tc.expCoins, senderAddrCoins)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOnTimeoutPacket() {
	baseDenom := "stake"
	senderAddr := sdk.AccAddress([]byte{0x1})
	receiveAddr := sdk.AccAddress([]byte{0x2})
	transferAmount := sdk.NewInt(100)
	testCases := []struct {
		name         string
		malleate     func(pxIbcTransferMsg *channeltypes.Packet)
		expPass      bool
		errorStr     string
		checkBalance bool
		expCoins     sdk.Coins
	}{
		{
			"pass - normal - ibc transfer packet",
			func(packet *channeltypes.Packet) {
				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				mintCoin(suite.T(), suite.chainPx.GetContext(), suite.GetApp(suite.chainPx.App), escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)))
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)),
		},
		{
			"pass - normal - px ibc transfer packet",
			func(packet *channeltypes.Packet) {
				packetData := pxtransfertypes.FungibleTokenPacketData{}
				pxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Router = ""
				packetData.Fee = sdk.ZeroInt().String()
				packet.Data = packetData.GetBytes()

				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				mintCoin(suite.T(), suite.chainPx.GetContext(), suite.GetApp(suite.chainPx.App), escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)))
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)),
		},
		{
			"pass - normal - ibc mint token - router is empty",
			func(packet *channeltypes.Packet) {
				packetData := pxtransfertypes.FungibleTokenPacketData{}
				pxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Denom = pxtypes.MintDenom()
				packetData.Router = ""
				packetData.Fee = sdk.ZeroInt().String()
				packet.Data = packetData.GetBytes()
				// mint coin to channel escrowAddress

				amount := sdk.NewCoins(sdk.NewCoin(pxtypes.MintDenom(), transferAmount))
				escrowAddress := transfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
				mintCoin(suite.T(), suite.chainPx.GetContext(), suite.GetApp(suite.chainPx.App), escrowAddress, amount)
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(pxtypes.MintDenom(), transferAmount)),
		},
		{
			"pass - router not empty | amount + zero fee",
			func(packet *channeltypes.Packet) {
				packetData := pxtransfertypes.FungibleTokenPacketData{}
				pxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Router = rand.Str(4)
				packetData.Fee = sdk.ZeroInt().String()
				packet.Data = packetData.GetBytes()

				amount := sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount))
				escrowAddress := transfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
				mintCoin(suite.T(), suite.chainPx.GetContext(), suite.GetApp(suite.chainPx.App), escrowAddress, amount)
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)),
		},
		{
			"pass - router not empty | amount + fee",
			func(packet *channeltypes.Packet) {
				packetData := pxtransfertypes.FungibleTokenPacketData{}
				pxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Router = rand.Str(4)
				fee := sdk.NewInt(50)
				packetData.Fee = fee.String()
				packet.Data = packetData.GetBytes()

				amount := sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.Add(fee)))
				escrowAddress := transfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
				mintCoin(suite.T(), suite.chainPx.GetContext(), suite.GetApp(suite.chainPx.App), escrowAddress, amount)
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.Add(sdk.NewInt(50)))),
		},
		{
			"error - escrow address insufficient 10coin",
			func(packet *channeltypes.Packet) {
				packetData := pxtransfertypes.FungibleTokenPacketData{}
				pxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Router = ""
				packetData.Fee = sdk.ZeroInt().String()
				packet.Data = packetData.GetBytes()

				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				mintCoin(suite.T(), suite.chainPx.GetContext(), suite.GetApp(suite.chainPx.App), escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.Sub(sdk.NewInt(10)))))
			},
			false,
			fmt.Sprintf("unable to unescrow tokens, this may be caused by a malicious counterparty module or a bug: please open an issue on counterparty module: %d%s is smaller than %d%s: insufficient funds", transferAmount.Sub(sdk.NewInt(10)).Uint64(), baseDenom, transferAmount.Uint64(), baseDenom),
			true,
			sdk.NewCoins(),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			chain := suite.GetApp(suite.chainPx.App)
			transferIBCModule := transfer.NewIBCModule(chain.TransferKeeper)
			pundixIBCMiddleware := pundixtransfer.NewIBCModule(chain.PundixTransferKeeper, transferIBCModule)
			packetData := transfertypes.NewFungibleTokenPacketData(baseDenom, transferAmount.String(), senderAddr.String(), receiveAddr.String())
			// only use timeout height
			packet := channeltypes.NewPacket(packetData.GetBytes(), 1, ibcgotesting.TransferPort, "channel-0", ibcgotesting.TransferPort, "channel-0", clienttypes.Height{
				RevisionNumber: 100,
				RevisionHeight: 100000,
			}, 0)
			tc.malleate(&packet)

			err := pundixIBCMiddleware.OnTimeoutPacket(suite.chainPx.GetContext(), packet, nil)
			if tc.expPass {
				suite.Require().NoError(err, "packetData:%s", string(packet.GetData()))
			} else {
				suite.Require().NotNil(err)
				suite.Require().Equalf(tc.errorStr, err.Error(), "packetData:%s", string(packet.GetData()))
			}

			if tc.checkBalance {
				bankKeeper := suite.GetApp(suite.chainPx.App).BankKeeper
				senderAddrCoins := bankKeeper.GetAllBalances(suite.chainPx.GetContext(), senderAddr)
				suite.Require().True(tc.expCoins.IsEqual(senderAddrCoins), "exp:%s,actual:%s", tc.expCoins, senderAddrCoins)
			}
		})
	}
}

func mintCoin(t *testing.T, ctx sdk.Context, chain *app.PundixApp, address sdk.AccAddress, coins sdk.Coins) {
	require.NoError(t, chain.BankKeeper.MintCoins(ctx, transfertypes.ModuleName, coins))
	require.NoError(t, chain.BankKeeper.SendCoinsFromModuleToAccount(ctx, transfertypes.ModuleName, address, coins))
}
