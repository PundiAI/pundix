package types

import (
	"fmt"
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
)

func TestParamsValidate(t *testing.T) {
	tests := []struct {
		name    string
		params  minttypes.Params
		expPass bool
	}{
		{name: "pass - default params", params: minttypes.DefaultParams(), expPass: true},

		{name: "error - empty mint denom", params: minttypes.NewParams("", sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.OneDec(), 1), expPass: false},
		{name: "error - invalid mint denom", params: minttypes.NewParams("1213123_2312", sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.OneDec(), 1), expPass: false},
		{name: "error - negative InflationRateChange", params: minttypes.NewParams("stake", sdk.NewDec(-1), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), 0), expPass: false},
		{name: "error - negative InflationMax", params: minttypes.NewParams("stake", sdk.ZeroDec(), sdk.NewDec(-1), sdk.ZeroDec(), sdk.ZeroDec(), 0), expPass: false},
		{name: "error - negative InflationMin", params: minttypes.NewParams("stake", sdk.ZeroDec(), sdk.ZeroDec(), sdk.NewDec(-1), sdk.ZeroDec(), 0), expPass: false},
		{name: "error - negative GoalBonded", params: minttypes.NewParams("stake", sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.NewDec(-1), 0), expPass: false},
		{name: "error - zero blocksPerYear", params: minttypes.NewParams("stake", sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.OneDec(), 0), expPass: false},
		{name: "error - InflationMin grant InflationMax", params: minttypes.NewParams("stake", sdk.ZeroDec(), sdk.NewDec(1), sdk.NewDec(2), sdk.OneDec(), 1), expPass: false},

		{name: "error - GoalBonded more then one", params: minttypes.NewParams("stake", sdk.OneDec(), sdk.OneDec(), sdk.OneDec(), sdk.NewDecWithPrec(11, 1), 1), expPass: false},
		{name: "pass - InflationRateChange more then one", params: minttypes.NewParams("stake", sdk.OneDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.OneDec(), 1), expPass: true},
		{name: "pass - InflationRateChange InflationMax InflationMin grate one", params: minttypes.NewParams("stake", sdk.NewDec(40), sdk.NewDec(40), sdk.NewDec(20), sdk.NewDecWithPrec(51, 2), minttypes.DefaultParams().BlocksPerYear), expPass: true},
	}
	require.NotPanics(t, func() {
		_ = ParamKeyTable()
	})
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Case %s", tt.name), func(t *testing.T) {
			genesisState := minttypes.NewGenesisState(minttypes.DefaultInitialMinter(), tt.params)

			validateGenesisErr := ValidateGenesis(*genesisState)
			if tt.expPass {
				require.NoError(t, validateGenesisErr)
			} else {
				require.Error(t, validateGenesisErr)
			}

			params := Params{Params: tt.params}
			validateParamsErr := params.Validate()
			if tt.expPass {
				require.NoError(t, validateParamsErr)
				for _, pair := range params.ParamSetPairs() {
					v := reflect.Indirect(reflect.ValueOf(pair.Value)).Interface()
					require.NoError(t, pair.ValidatorFn(v))
				}
			} else {
				require.Error(t, validateParamsErr)
			}
		})
	}
}

func TestValidateMintDenom(t *testing.T) {
	require.NoError(t, validateMintDenom("stake"))
	require.Error(t, validateMintDenom(""))
	require.Error(t, validateMintDenom(1111))
	require.Error(t, validateMintDenom("--11--"))
}

func TestValidateInflationRateChange(t *testing.T) {
	require.NoError(t, validateInflationRateChange(sdk.OneDec()))
	require.Errorf(t, validateInflationRateChange("1"), "error - type error")
	require.Errorf(t, validateInflationRateChange(sdk.NewDec(-1)), "error - negative")
}

func TestValidateInflationMax(t *testing.T) {
	require.NoError(t, validateInflationMax(sdk.OneDec()))
	require.Errorf(t, validateInflationMax("1"), "error - type error")
	require.Errorf(t, validateInflationMax(sdk.NewDec(-1)), "error - negative")
}

func TestValidateInflationMin(t *testing.T) {
	require.NoError(t, validateInflationMin(sdk.OneDec()))
	require.Errorf(t, validateInflationMin("1"), "error - type error")
	require.Errorf(t, validateInflationMin(sdk.NewDec(-1)), "error - negative")
}

func TestValidateGoalBonded(t *testing.T) {
	require.NoError(t, validateGoalBonded(sdk.OneDec()))
	require.Errorf(t, validateGoalBonded(""), "error - type error")
	require.Errorf(t, validateGoalBonded(sdk.NewDec(-1)), "error - negative")
	require.Errorf(t, validateGoalBonded(sdk.ZeroDec()), "error - zero")
	require.Errorf(t, validateGoalBonded(sdk.NewDec(2)), "error - grate than one")
}

func TestValidateBlocksPerYear(t *testing.T) {
	require.NoError(t, validateBlocksPerYear(uint64(1)))
	require.Errorf(t, validateBlocksPerYear(uint64(0)), "error - zero value")
	require.Errorf(t, validateBlocksPerYear("1"), "error - type error")
}
