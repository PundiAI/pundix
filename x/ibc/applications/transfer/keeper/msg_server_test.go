package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	pxibctesting "github.com/pundix/pundix/ibc/testing"
	pxtransfertypes "github.com/pundix/pundix/x/ibc/applications/transfer/types"
)

func (suite *KeeperTestSuite) TestTransfer() {

	var channel0IbcDenom = transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: "stake",
	}
	testCases := []struct {
		name     string
		malleate func(pxIbcTransferMsg *pxtransfertypes.MsgTransfer, ibcTransferMsg *transfertypes.MsgTransfer)
		expPass  bool
		errorStr string
	}{
		{
			"pass - normal",
			func(_ *pxtransfertypes.MsgTransfer, _ *transfertypes.MsgTransfer) {
			},
			true,
			"",
		},
		{
			"pass - normal - router + fee",
			func(pxIbcTransferMsg *pxtransfertypes.MsgTransfer, ibcTransferMsg *transfertypes.MsgTransfer) {
				fee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1))
				pxIbcTransferMsg.Fee = fee
				pxIbcTransferMsg.Router = "cosmos"
				suite.Require().NoError(suite.GetApp(suite.chainPx.App).BankKeeper.MintCoins(suite.chainPx.GetContext(), transfertypes.ModuleName, sdk.NewCoins(fee)))
				suite.Require().NoError(suite.GetApp(suite.chainPx.App).BankKeeper.SendCoinsFromModuleToAccount(suite.chainPx.GetContext(), transfertypes.ModuleName, suite.chainPx.SenderAccount.GetAddress(), sdk.NewCoins(fee)))
				ibcTransferMsg.Sender = ""
			},
			true,
			"",
		},
		{
			"error - invalid sender",
			func(pxIbcTransferMsg *pxtransfertypes.MsgTransfer, ibcTransferMsg *transfertypes.MsgTransfer) {
				ibcTransferMsg.Sender = "address"
				pxIbcTransferMsg.Sender = "address"
			},
			false,
			"decoding bech32 failed: invalid bech32 string length 7",
		},
		{
			"error - send is disable",
			func(pxIbcTransferMsg *pxtransfertypes.MsgTransfer, ibcTransferMsg *transfertypes.MsgTransfer) {
				suite.GetApp(suite.chainPx.App).TransferKeeper.SetParams(suite.chainPx.GetContext(), transfertypes.NewParams(false, true))
			},
			false,
			"fungible token transfers from this chain are disabled",
		},
		{
			name: "error - coin send is disable",
			malleate: func(pxIbcTransferMsg *pxtransfertypes.MsgTransfer, ibcTransferMsg *transfertypes.MsgTransfer) {
				suite.GetApp(suite.chainPx.App).BankKeeper.SetParams(suite.chainPx.GetContext(), banktypes.NewParams(true, []*banktypes.SendEnabled{{Denom: sdk.DefaultBondDenom, Enabled: false}}))
			},
			errorStr: fmt.Sprintf("%s transfers are currently disabled: fungible token transfers from this chain are disabled", "stake"),
		},
		{
			"error - sender is blocked addr",
			func(pxIbcTransferMsg *pxtransfertypes.MsgTransfer, ibcTransferMsg *transfertypes.MsgTransfer) {
				ibcTransferMsg.Sender = authtypes.NewModuleAddress(transfertypes.ModuleName).String()
				pxIbcTransferMsg.Sender = authtypes.NewModuleAddress(transfertypes.ModuleName).String()
			},
			false,
			fmt.Sprintf("%s is not allowed to send funds: unauthorized", authtypes.NewModuleAddress(transfertypes.ModuleName).String()),
		},
		{
			"error - channel not found",
			func(pxIbcTransferMsg *pxtransfertypes.MsgTransfer, ibcTransferMsg *transfertypes.MsgTransfer) {
				ibcTransferMsg.SourceChannel = "channel-1"
				pxIbcTransferMsg.SourceChannel = "channel-1"
			},
			false,
			fmt.Sprintf("port ID (%s) channel ID (%s): channel not found", "transfer", "channel-1"),
		},
		{
			"err - sender balance insufficient - not route - only amount",
			func(pxIbcTransferMsg *pxtransfertypes.MsgTransfer, ibcTransferMsg *transfertypes.MsgTransfer) {
				amount := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))
				pxIbcTransferMsg.Sender = sdk.AccAddress([]byte{0x1}).String()
				pxIbcTransferMsg.Token = amount

				ibcTransferMsg.Sender = sdk.AccAddress([]byte{0x1}).String()
				ibcTransferMsg.Token = amount

			},
			false,
			fmt.Sprintf("%s is smaller than %s: insufficient funds", "0stake", "10stake"),
		},
		{
			"err - sender balance insufficient - has route - amount + fee",
			func(pxIbcTransferMsg *pxtransfertypes.MsgTransfer, ibcTransferMsg *transfertypes.MsgTransfer) {
				amount := sdk.NewCoin("stake", sdk.NewInt(10))
				fee := sdk.NewCoin("stake", sdk.NewInt(10))
				pxIbcTransferMsg.Sender = sdk.AccAddress([]byte{0x1}).String()
				pxIbcTransferMsg.Token = amount
				pxIbcTransferMsg.Router = "cosmos"
				pxIbcTransferMsg.Fee = fee

				ibcTransferMsg.Sender = ""
			},
			false,
			fmt.Sprintf("%s is smaller than %s: insufficient funds", "0stake", "20stake"),
		},
		{
			"error - ibc denom not found",
			func(pxIbcTransferMsg *pxtransfertypes.MsgTransfer, ibcTransferMsg *transfertypes.MsgTransfer) {
				amount := sdk.NewCoin(channel0IbcDenom.IBCDenom(), sdk.NewInt(10))
				pxIbcTransferMsg.Token = amount
				ibcTransferMsg.Token = amount
			},
			false,
			fmt.Sprintf("%s: denomination trace not found", channel0IbcDenom.Hash()),
		},
		{
			"error - sender ibc coin insufficient funds - not include fee",
			func(pxIbcTransferMsg *pxtransfertypes.MsgTransfer, ibcTransferMsg *transfertypes.MsgTransfer) {
				suite.GetApp(suite.chainPx.App).TransferKeeper.SetDenomTrace(suite.chainPx.GetContext(), channel0IbcDenom)
				amount := sdk.NewCoin(channel0IbcDenom.IBCDenom(), sdk.NewInt(10))
				pxIbcTransferMsg.Token = amount
				ibcTransferMsg.Token = amount
			},
			false,
			fmt.Sprintf("%s%s is smaller than %s%s: insufficient funds", "0", channel0IbcDenom.IBCDenom(), "10", channel0IbcDenom.IBCDenom()),
		},
	}

	var covertPxIBCToCosmosIBCTransfer = func(msg *pxtransfertypes.MsgTransfer) *transfertypes.MsgTransfer {
		result := transfertypes.NewMsgTransfer(msg.SourcePort, msg.SourceChannel, msg.Token, msg.Sender, msg.Receiver, msg.TimeoutHeight, msg.TimeoutTimestamp)
		result.Memo = msg.Memo
		return result
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			// init chain
			suite.SetupTest()
			path := pxibctesting.NewTransferPath(suite.chainPx, suite.chainFx)
			// init channel
			suite.coordinator.Setup(path)

			// mint transfer coin to sender account
			coin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))
			fee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0))
			thisChainApp := suite.GetApp(suite.chainPx.App)
			suite.Require().NoError(thisChainApp.BankKeeper.MintCoins(suite.chainPx.GetContext(), transfertypes.ModuleName, sdk.NewCoins(coin)))
			suite.Require().NoError(thisChainApp.BankKeeper.SendCoinsFromModuleToAccount(suite.chainPx.GetContext(), transfertypes.ModuleName, suite.chainPx.SenderAccount.GetAddress(), sdk.NewCoins(coin)))

			// build pxIBCTransfer
			pxIbcTransferMsg := pxtransfertypes.NewMsgTransfer(
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
				coin,
				suite.chainPx.SenderAccount.GetAddress().String(),
				suite.chainFx.SenderAccount.GetAddress().String(),
				suite.chainFx.GetTimeoutHeight(),
				0, // only use timeout height
				"",
				fee,
			)
			pxIbcTransferMsg.Memo = "memo"

			// build cosmosIBCTransfer
			ibcTransferMsg := covertPxIBCToCosmosIBCTransfer(pxIbcTransferMsg)

			// test case maybe update transferMsg attribute
			tc.malleate(pxIbcTransferMsg, ibcTransferMsg)

			// test:1 call pundix ibc transfer
			cacheContext, _ := suite.chainPx.GetContext().CacheContext()
			_, err := thisChainApp.PundixTransferKeeper.Transfer(sdk.WrapSDKContext(cacheContext), pxIbcTransferMsg)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(err.Error(), tc.errorStr)
			}

			// if ibc transfer msg is nil or sender is empty, return: only test px ibcTransferMsg
			if ibcTransferMsg == nil || len(ibcTransferMsg.Sender) == 0 {
				return
			}

			// test:2 call cosmos ibc transfer
			cacheContext, _ = suite.chainPx.GetContext().CacheContext()
			_, err = thisChainApp.TransferKeeper.Transfer(sdk.WrapSDKContext(cacheContext), ibcTransferMsg)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(err.Error(), tc.errorStr)
			}
		})
	}
}
