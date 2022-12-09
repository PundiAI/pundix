package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcgotesting "github.com/cosmos/ibc-go/v3/testing"
	pxtypes "github.com/pundix/pundix/types"
	pundixtransfer "github.com/pundix/pundix/x/ibc/applications/transfer"
	pxtransfertypes "github.com/pundix/pundix/x/ibc/applications/transfer/types"
	"github.com/tendermint/tendermint/libs/rand"
)

func (suite *KeeperTestSuite) TestOnRecvPacket() {
	baseDenom := "stake"
	senderAddr := sdk.AccAddress([]byte{0x1}).String()
	receiveAddr := sdk.AccAddress([]byte{0x2}).String()
	ibcDenomTrace := transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: pxtypes.MintDenom(),
	}
	ibcAmount := sdk.NewInt(100)
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
			sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), ibcAmount)),
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
			sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), ibcAmount)),
		},
		{
			"pass - normal - ibc purse token - router is empty",
			func(packet *channeltypes.Packet) {
				packetData := pxtransfertypes.FungibleTokenPacketData{}
				pxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Denom = ibcDenomTrace.IBCDenom()
				packetData.Router = ""
				packetData.Fee = sdk.ZeroInt().String()
				packet.Data = packetData.GetBytes()
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), ibcAmount)),
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
			"not support router",
			true,
			sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), ibcAmount)),
		},
		{
			"error - not support fee",
			func(packet *channeltypes.Packet) {
				packetData := pxtransfertypes.FungibleTokenPacketData{}
				pxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Denom = ibcDenomTrace.IBCDenom()
				packetData.Router = rand.Str(4)
				packetData.Fee = sdk.NewInt(rand.Int64()).String()
				packet.Data = packetData.GetBytes()
			},
			false,
			"not support router",
			true,
			sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), ibcAmount)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			app := suite.GetApp(suite.chainPx.App)
			transferIBCModule := transfer.NewIBCModule(app.TransferKeeper)
			pundixIBCMiddleware := pundixtransfer.NewIBCModule(app.PundixTransferKeeper, transferIBCModule)
			packetData := transfertypes.NewFungibleTokenPacketData(baseDenom, ibcAmount.String(), senderAddr, receiveAddr)
			// only use timeout height
			packet := channeltypes.NewPacket(packetData.GetBytes(), 1, ibcgotesting.TransferPort, "channel-0", ibcgotesting.TransferPort, "channel-0", clienttypes.Height{
				RevisionNumber: 100,
				RevisionHeight: 100000,
			}, 0)
			tc.malleate(&packet)

			ack := pundixIBCMiddleware.OnRecvPacket(suite.chainPx.GetContext(), packet, nil)
			suite.Require().NotNil(ack)
			if tc.expPass {
				suite.Require().True(ack.Success())
			} else {
				suite.Require().False(ack.Success())
				_, ok := ack.(channeltypes.Acknowledgement)
				suite.Require().True(ok)
				// suite.Require().Equal(tc.errorStr, acknowledgement.GetError())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOnAcknowledgementPacket() {
}

func (suite *KeeperTestSuite) TestOnTimeoutPacket() {
}
