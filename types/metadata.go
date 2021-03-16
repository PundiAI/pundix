package types

import (
	"strings"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func GetPURSEMetaData(denom string) banktypes.Metadata {
	return banktypes.Metadata{
		Description: "PURSE Token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    strings.ToLower(denom),
				Exponent: 0,
				Aliases:  nil,
			},
			{
				Denom:    "PURSE",
				Exponent: 18,
				Aliases:  nil,
			},
		},
		Base:    strings.ToLower(denom),
		Name:    "PURSE Token",
		Symbol:  "PURSE",
		Display: "PURSE",
	}
}

func GetPUNDIXMetaData(denom string) banktypes.Metadata {
	return banktypes.Metadata{
		Description: "Pundi X Token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    strings.ToLower(denom),
				Exponent: 0,
				Aliases:  nil,
			},
			{
				Denom:    "PUNDIX",
				Exponent: 18,
				Aliases:  nil,
			},
		},
		Base:    strings.ToLower(denom),
		Name:    "Pundi X Token",
		Symbol:  "PUNDIX",
		Display: "PUNDIX",
	}
}
