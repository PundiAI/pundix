package types_test

import (
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	types "github.com/pundix/pundix/x/ibc/applications/transfer/types"

	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

const (
	denom              = "transfer/gaiachannel/atom"
	amount             = "100"
	largeAmount        = "18446744073709551616"                                                           // one greater than largest uint64 (^uint64(0))
	invalidLargeAmount = "115792089237316195423570985008687907853269984665640564039457584007913129639936" // 2^256
)

var (
	defaultPacketRouter = ""
	defaultPacketFeeStr = sdk.ZeroInt().String()
)

// TestFungibleTokenPacketDataValidateBasic tests ValidateBasic for FungibleTokenPacketData
func TestFungibleTokenPacketDataValidateBasic(t *testing.T) {
	testCases := []struct {
		name       string
		packetData types.FungibleTokenPacketData
		expPass    bool
	}{
		{"valid packet", types.NewFungibleTokenPacketData(denom, amount, addr1, addr2, defaultPacketRouter, defaultPacketFeeStr), true},
		{"valid packet with large amount", types.NewFungibleTokenPacketData(denom, largeAmount, addr1, addr2, defaultPacketRouter, defaultPacketFeeStr), true},
		{"invalid denom", types.NewFungibleTokenPacketData("", amount, addr1, addr2, defaultPacketRouter, defaultPacketFeeStr), false},
		{"invalid empty amount", types.NewFungibleTokenPacketData(denom, "", addr1, addr2, defaultPacketRouter, defaultPacketFeeStr), false},
		{"invalid zero amount", types.NewFungibleTokenPacketData(denom, "0", addr1, addr2, defaultPacketRouter, defaultPacketFeeStr), false},
		{"invalid negative amount", types.NewFungibleTokenPacketData(denom, "-1", addr1, addr2, defaultPacketRouter, defaultPacketFeeStr), false},
		{"invalid large amount", types.NewFungibleTokenPacketData(denom, invalidLargeAmount, addr1, addr2, defaultPacketRouter, defaultPacketFeeStr), false},
		{"missing sender address", types.NewFungibleTokenPacketData(denom, amount, emptyAddr, addr2, defaultPacketRouter, defaultPacketFeeStr), false},
		{"missing recipient address", types.NewFungibleTokenPacketData(denom, amount, addr1, emptyAddr, defaultPacketRouter, defaultPacketFeeStr), false},
		{"invalid negative fee", types.NewFungibleTokenPacketData(denom, amount, addr1, emptyAddr, defaultPacketRouter, "-1"), false},
	}

	for i, tc := range testCases {
		err := tc.packetData.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "valid test case %d failed: %v", i, err)
		} else {
			require.Error(t, err, "invalid test case %d passed: %s", i, tc.name)
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name   string
		data   []byte
		pass   bool
		expErr error
		exp    types.FungibleTokenPacketData
	}{
		{
			name:   "fx transfer packet - no router",
			data:   types.NewFungibleTokenPacketData("FX", "100", "Add1", "Add2", "", "0").GetBytes(),
			pass:   true,
			expErr: nil,
			exp: types.FungibleTokenPacketData{
				Denom:    "FX",
				Amount:   "100",
				Sender:   "Add1",
				Receiver: "Add2",
				Router:   "",
				Fee:      "",
			},
		},
		{
			name:   "fx transfer packet - router with 0fee",
			data:   types.NewFungibleTokenPacketData("FX", "100", "Add1", "Add2", "router", "0").GetBytes(),
			pass:   true,
			expErr: nil,
			exp: types.FungibleTokenPacketData{
				Denom:    "FX",
				Amount:   "100",
				Sender:   "Add1",
				Receiver: "Add2",
				Router:   "router",
				Fee:      "0",
			},
		},
		{
			name:   "fx transfer packet - router with empty fee",
			data:   types.NewFungibleTokenPacketData("FX", "100", "Add1", "Add2", "router", "").GetBytes(),
			pass:   true,
			expErr: nil,
			exp: types.FungibleTokenPacketData{
				Denom:    "FX",
				Amount:   "100",
				Sender:   "Add1",
				Receiver: "Add2",
				Router:   "router",
				Fee:      "",
			},
		},
		{
			name:   "ibc transfer packet",
			data:   transfertypes.NewFungibleTokenPacketData("FX", "100", "Add1", "Add2").GetBytes(),
			pass:   true,
			expErr: nil,
			exp: types.FungibleTokenPacketData{
				Denom:    "FX",
				Amount:   "100",
				Sender:   "Add1",
				Receiver: "Add2",
				Router:   "",
				Fee:      "",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var packet types.FungibleTokenPacketData
			err := types.ModuleCdc.UnmarshalJSON(testCase.data, &packet)
			if testCase.pass {
				require.NoError(t, err)
				require.EqualValues(t, testCase.exp, packet)
				require.EqualValues(t, testCase.exp.GetDenom(), packet.GetDenom())
				require.EqualValues(t, testCase.exp.GetAmount(), packet.GetAmount())
				require.EqualValues(t, testCase.exp.GetSender(), packet.GetSender())
				require.EqualValues(t, testCase.exp.GetReceiver(), packet.GetReceiver())
				require.EqualValues(t, testCase.exp.GetRouter(), packet.GetRouter())
				require.EqualValues(t, testCase.exp.GetFee(), packet.GetFee())
			} else {
				require.ErrorIs(t, err, testCase.expErr)
			}
		})
	}
}
