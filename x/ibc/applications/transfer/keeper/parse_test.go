package keeper

import (
	"bytes"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/pundix/pundix/x/ibc/applications/transfer/types"
)

func TestParseReceiveAndAmountByPacket(t *testing.T) {
	type expect struct {
		address []byte
		amount  sdk.Int
		fee     sdk.Int
	}
	testCases := []struct {
		name    string
		packet  types.FungibleTokenPacketData
		expPass bool
		err     error
		expect  expect
	}{
		{
			"no router - expect address is receive",
			types.FungibleTokenPacketData{Receiver: sdk.AccAddress("receive1").String(), Amount: "1", Fee: "0"},
			true, nil,
			expect{address: sdk.AccAddress("receive1"), amount: sdk.NewIntFromUint64(1), fee: sdk.NewIntFromUint64(0)},
		},
		{
			"no router - expect fee is 0, input 1",
			types.FungibleTokenPacketData{Receiver: sdk.AccAddress("receive1").String(), Amount: "1", Fee: "0"},
			true, nil,
			expect{address: sdk.AccAddress("receive1"), amount: sdk.NewIntFromUint64(1), fee: sdk.NewIntFromUint64(0)},
		},
		{
			"router - expect address is sender",
			types.FungibleTokenPacketData{Sender: sdk.AccAddress("sender1").String(), Receiver: sdk.AccAddress("receive1").String(), Amount: "1", Fee: "0", Router: "erc20"},
			true, nil,
			expect{address: sdk.AccAddress("sender1"), amount: sdk.NewIntFromUint64(1), fee: sdk.NewIntFromUint64(0)},
		},
		{
			"router - expect fee is 1, input 1",
			types.FungibleTokenPacketData{Sender: sdk.AccAddress("sender1").String(), Receiver: sdk.AccAddress("receive1").String(), Amount: "1", Fee: "1", Router: "erc20"},
			true, nil,
			expect{address: sdk.AccAddress("sender1"), amount: sdk.NewIntFromUint64(1), fee: sdk.NewIntFromUint64(1)},
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualAddress, actualAmount, actualFee, err := parseReceiveAndAmountByPacket(tc.packet)
			if tc.expPass {
				require.NoError(t, err, "valid test case %d failed: %v", i, err)
			} else {
				require.Error(t, err)
				require.EqualValues(t, err.Error(), tc.err.Error())
			}
			require.Truef(t, bytes.Equal(tc.expect.address, actualAddress.Bytes()), "expected %s, actual %s", sdk.AccAddress(tc.expect.address).String(), actualAddress.String())
			require.EqualValues(t, tc.expect.amount.String(), actualAmount.String())
			require.EqualValues(t, tc.expect.fee.String(), actualFee.String())
		})
	}
}

func TestParseAmountAndFeeByPacket(t *testing.T) {
	type expect struct {
		amount sdk.Int
		fee    sdk.Int
	}
	testCases := []struct {
		name    string
		packet  types.FungibleTokenPacketData
		expPass bool
		errStr  string
		expect  expect
	}{
		{
			"pass - no router only amount ",
			types.FungibleTokenPacketData{Amount: "1"},
			true, "",
			expect{amount: sdk.NewInt(1), fee: sdk.ZeroInt()},
		},
		{
			"error - amount is empty",
			types.FungibleTokenPacketData{Amount: ""},
			false,
			"unable to parse transfer amount () into sdk.Int: invalid token amount",
			expect{amount: sdk.Int{}, fee: sdk.Int{}},
		},
		{
			"error - fee is empty",
			types.FungibleTokenPacketData{Amount: "1", Fee: "", Router: "aaa"},
			false,
			"fee amount is invalid:: invalid token amount",
			expect{amount: sdk.Int{}, fee: sdk.Int{}},
		},
		{
			"error - fee is negative",
			types.FungibleTokenPacketData{Amount: "1", Fee: "-1", Router: "aaa"},
			false,
			"fee amount is invalid:-1: invalid token amount",
			expect{amount: sdk.Int{}, fee: sdk.Int{}},
		},
		{
			"pass - fee is zero",
			types.FungibleTokenPacketData{Amount: "1", Fee: "0", Router: "aaa"},
			true,
			"",
			expect{amount: sdk.NewInt(1), fee: sdk.ZeroInt()},
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualAmount, actualFee, err := parseAmountAndFeeByPacket(tc.packet)
			if tc.expPass {
				require.NoError(t, err, "valid test case %d failed: %v", i, err)
			} else {
				require.Error(t, err)
				require.EqualValues(t, tc.errStr, err.Error())
			}
			require.EqualValues(t, tc.expect.amount.String(), actualAmount.String())
			require.EqualValues(t, tc.expect.fee.String(), actualFee.String())
		})
	}
}

func TestParsePacketAddress(t *testing.T) {
	testCases := []struct {
		name    string
		address string
		expPass bool
		err     error
		expect  []byte
	}{
		{"normal fx address", sdk.AccAddress("abc").String(), true, nil, sdk.AccAddress("abc")},

		{"err bech32 address - kc74", "fx1yef9232palu3ps25ldjr62ck046rgd8292kc74", false, fmt.Errorf("decoding bech32 failed: invalid checksum (expected 92kc73 got 92kc74)"), []byte{}},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualAddress, err := parsePacketAddress(tc.address)
			if tc.expPass {
				require.NoError(t, err, "valid test case %d failed: %v", i, err)
			} else {
				require.Error(t, err)
				require.EqualValues(t, err.Error(), tc.err.Error())
			}
			require.Truef(t, bytes.Equal(tc.expect, actualAddress.Bytes()), "expected %s, actual %s", sdk.AccAddress(tc.expect).String(), actualAddress.String())
		})
	}
}
