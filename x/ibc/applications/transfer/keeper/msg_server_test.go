package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/stretchr/testify/mock"

	"github.com/pundix/pundix/x/ibc/applications/transfer/keeper"
	pxtransfertypes "github.com/pundix/pundix/x/ibc/applications/transfer/types"
)

func (suite *KeeperTestSuite) TestTransfer() {
	mockChannelKeeper := &MockChannelKeeper{}
	mockICS4Wrapper := &MockICS4Wrapper{}
	mockChannelKeeper.On("GetNextSequenceSend", mock.Anything, mock.Anything, "channel-0").Return(uint64(1), true)
	mockChannelKeeper.On("GetNextSequenceSend", mock.Anything, mock.Anything, "channel-2").Return(uint64(0), false)
	mockChannelKeeper.On("GetNextSequenceSend", mock.Anything, mock.Anything, "channel-3").Return(uint64(1), true)
	mockChannelKeeper.On("GetChannel", mock.Anything, mock.Anything, "channel-0").Return(channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-1")}, true)
	mockChannelKeeper.On("GetChannel", mock.Anything, mock.Anything, "channel-1").Return(channeltypes.Channel{}, false)
	mockChannelKeeper.On("GetChannel", mock.Anything, mock.Anything, "channel-2").Return(channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-3")}, true)
	mockChannelKeeper.On("GetChannel", mock.Anything, mock.Anything, "channel-3").Return(channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-4")}, true)
	mockICS4Wrapper.On("SendPacket", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	channel1IbcDenom := types.DenomTrace{
		Path:      "transfer/channel-1",
		BaseDenom: "stake",
	}
	testCases := []struct {
		name     string
		malleate func() *pxtransfertypes.MsgTransfer
		expPass  bool
		errorStr string
	}{
		{
			"pass - normal",
			func() *pxtransfertypes.MsgTransfer {
				senderAcc := sdk.AccAddress([]byte{0x001})
				transferMsg := pxtransfertypes.NewMsgTransfer("transfer", "channel-0", sdk.NewCoin("stake", sdk.NewInt(10)), senderAcc.String(), "", timeoutHeight, 0, "", sdk.NewCoin("stake", sdk.NewInt(0)))

				coins := sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(10)))
				err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, senderAcc, coins)
				suite.Require().NoError(err)
				suite.Commit()
				return transferMsg
			},
			true,
			"",
		},
		{
			"pass - normal - router + fee",
			func() *pxtransfertypes.MsgTransfer {
				senderAcc := sdk.AccAddress([]byte{0x001})
				amount := sdk.NewCoin("stake", sdk.NewInt(10))
				fee := sdk.NewCoin("stake", sdk.NewInt(10))
				transferMsg := pxtransfertypes.NewMsgTransfer("transfer", "channel-0", amount, senderAcc.String(), "", timeoutHeight, 0, "bsc", fee)

				coins := sdk.NewCoins(amount.Add(fee))
				err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, senderAcc, coins)
				suite.Require().NoError(err)
				suite.Commit()
				return transferMsg
			},
			true,
			"",
		},
		{
			"error - invalid sender addr",
			func() *pxtransfertypes.MsgTransfer {
				transferMsg := pxtransfertypes.NewMsgTransfer("transfer", "channel-0", sdk.NewCoin("stake", sdk.NewInt(10)), "0xxxxx", "", timeoutHeight, 0, "", sdk.NewCoin("stake", sdk.NewInt(0)))
				return transferMsg
			},
			false,
			"decoding bech32 failed: invalid bech32 string length 6",
		},
		{
			"error - send is disable",
			func() *pxtransfertypes.MsgTransfer {
				params := suite.app.TransferKeeper.GetParams(suite.ctx)
				params.SendEnabled = false
				suite.app.TransferKeeper.SetParams(suite.ctx, params)
				suite.Commit()
				senderAcc := sdk.AccAddress([]byte{0x001})
				transferMsg := pxtransfertypes.NewMsgTransfer("transfer", "channel-0", sdk.NewCoin("stake", sdk.NewInt(10)), senderAcc.String(), "", timeoutHeight, 0, "", sdk.NewCoin("stake", sdk.NewInt(0)))
				return transferMsg
			},
			false,
			"fungible token transfers from this chain are disabled",
		},
		{
			"error - coin send is disable",
			func() *pxtransfertypes.MsgTransfer {
				bankParams := suite.app.BankKeeper.GetParams(suite.ctx)
				sendEnabled := bankParams.GetSendEnabled()
				if bankParams.SendEnabled == nil {
					bankParams.SendEnabled = []*banktypes.SendEnabled{}
				}
				found := false
				for _, tokenSend := range sendEnabled {
					if tokenSend.Denom == "stake" {
						tokenSend.Enabled = false
						found = true
						break
					}
				}
				if !found {
					sendEnabled = append(sendEnabled, banktypes.NewSendEnabled("stake", false))
				}
				bankParams.SendEnabled = sendEnabled
				suite.app.BankKeeper.SetParams(suite.ctx, bankParams)
				suite.Commit()
				senderAcc := sdk.AccAddress([]byte{0x001})
				transferMsg := pxtransfertypes.NewMsgTransfer("transfer", "channel-0", sdk.NewCoin("stake", sdk.NewInt(10)), senderAcc.String(), "", timeoutHeight, 0, "", sdk.NewCoin("stake", sdk.NewInt(0)))
				return transferMsg
			},
			false,
			fmt.Sprintf("%s transfers are currently disabled: fungible token transfers from this chain are disabled", "stake"),
		},
		{
			"error - sender is blocked addr",
			func() *pxtransfertypes.MsgTransfer {
				senderAcc := authtypes.NewModuleAddress(types.ModuleName)
				transferMsg := pxtransfertypes.NewMsgTransfer("transfer", "channel-0", sdk.NewCoin("stake", sdk.NewInt(10)), senderAcc.String(), "", timeoutHeight, 0, "", sdk.NewCoin("stake", sdk.NewInt(0)))
				return transferMsg
			},
			false,
			fmt.Sprintf("%s is not allowed to send funds: unauthorized", authtypes.NewModuleAddress(types.ModuleName).String()),
		},
		{
			"error - channel not found",
			func() *pxtransfertypes.MsgTransfer {
				senderAcc := sdk.AccAddress([]byte{0x001})
				transferMsg := pxtransfertypes.NewMsgTransfer("transfer", "channel-1", sdk.NewCoin("stake", sdk.NewInt(10)), senderAcc.String(), "", timeoutHeight, 0, "", sdk.NewCoin("stake", sdk.NewInt(0)))
				return transferMsg
			},
			false,
			fmt.Sprintf("port ID (%s) channel ID (%s): channel not found", "transfer", "channel-1"),
		},
		{
			"error - next sequence send not found",
			func() *pxtransfertypes.MsgTransfer {
				senderAcc := sdk.AccAddress([]byte{0x001})
				transferMsg := pxtransfertypes.NewMsgTransfer("transfer", "channel-2", sdk.NewCoin("stake", sdk.NewInt(10)), senderAcc.String(), "", timeoutHeight, 0, "", sdk.NewCoin("stake", sdk.NewInt(0)))
				return transferMsg
			},
			false,
			fmt.Sprintf("source port: transfer, source channel: %s: sequence send not found", "channel-2"),
		},
		{
			"error - capability not found",
			func() *pxtransfertypes.MsgTransfer {
				senderAcc := sdk.AccAddress([]byte{0x001})
				transferMsg := pxtransfertypes.NewMsgTransfer("transfer", "channel-3", sdk.NewCoin("stake", sdk.NewInt(10)), senderAcc.String(), "", timeoutHeight, 0, "", sdk.NewCoin("stake", sdk.NewInt(0)))
				return transferMsg
			},
			false,
			"module does not own channel capability: channel capability not found",
		},
		{
			"err - sender balance insufficient - not route - only amount",
			func() *pxtransfertypes.MsgTransfer {
				senderAcc := sdk.AccAddress([]byte{0x001})
				amount := sdk.NewCoin("stake", sdk.NewInt(10))
				fee := sdk.NewCoin("stake", sdk.NewInt(10))
				transferMsg := pxtransfertypes.NewMsgTransfer("transfer", "channel-0", amount, senderAcc.String(), "", timeoutHeight, 0, "", fee)
				return transferMsg
			},
			false,
			fmt.Sprintf("%s is smaller than %s: insufficient funds", "0stake", "10stake"),
		},
		{
			"err - sender balance insufficient - has route - amount + fee",
			func() *pxtransfertypes.MsgTransfer {
				senderAcc := sdk.AccAddress([]byte{0x001})
				amount := sdk.NewCoin("stake", sdk.NewInt(10))
				fee := sdk.NewCoin("stake", sdk.NewInt(10))
				transferMsg := pxtransfertypes.NewMsgTransfer("transfer", "channel-0", amount, senderAcc.String(), "", timeoutHeight, 0, "bsc", fee)
				return transferMsg
			},
			false,
			fmt.Sprintf("%s is smaller than %s: insufficient funds", "0stake", "20stake"),
		},
		{
			"error - ibc denom not found",
			func() *pxtransfertypes.MsgTransfer {
				senderAcc := sdk.AccAddress([]byte{0x001})
				transferMsg := pxtransfertypes.NewMsgTransfer("transfer", "channel-0", sdk.NewCoin(channel1IbcDenom.IBCDenom(), sdk.NewInt(10)), senderAcc.String(), "", timeoutHeight, 0, "", sdk.NewCoin("stake", sdk.NewInt(0)))
				return transferMsg
			},
			false,
			fmt.Sprintf("%s: denomination trace not found", channel1IbcDenom.Hash()),
		},
		{
			"error - sender ibc coin insufficient funds - not include fee",
			func() *pxtransfertypes.MsgTransfer {
				senderAcc := sdk.AccAddress([]byte{0x001})
				suite.app.TransferKeeper.SetDenomTrace(suite.ctx, channel1IbcDenom)
				suite.Commit()
				transferMsg := pxtransfertypes.NewMsgTransfer("transfer", "channel-0", sdk.NewCoin(channel1IbcDenom.IBCDenom(), sdk.NewInt(10)), senderAcc.String(), "", timeoutHeight, 0, "", sdk.NewCoin(channel1IbcDenom.IBCDenom(), sdk.NewInt(0)))
				return transferMsg
			},
			false,
			fmt.Sprintf("%s%s is smaller than %s%s: insufficient funds", "0", channel1IbcDenom.IBCDenom(), "10", channel1IbcDenom.IBCDenom()),
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()

			_, err := suite.app.ScopedTransferKeeper.NewCapability(suite.ctx, host.ChannelCapabilityPath("transfer", "channel-0"))
			suite.Require().NoError(err)

			msg := tc.malleate()

			// use mockChannelKeeper
			suite.app.PundixTransferKeeper = keeper.NewKeeper(suite.app.TransferKeeper,
				suite.app.AppCodec(), suite.app.GetKey(types.StoreKey), suite.app.GetSubspace(types.ModuleName),
				&MockICS4Wrapper{},
				mockChannelKeeper, &suite.app.IBCKeeper.PortKeeper,
				suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.ScopedTransferKeeper,
			)

			_, err = suite.app.PundixTransferKeeper.Transfer(sdk.WrapSDKContext(suite.ctx), msg)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(err.Error(), tc.errorStr)
			}
		})
	}
}
