package types

import minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data minttypes.GenesisState) error {
	if err := (Params{Params: data.Params}).Validate(); err != nil {
		return err
	}

	return minttypes.ValidateMinter(data.Minter)
}
