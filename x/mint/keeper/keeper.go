package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/pundix/pundix/x/mint/types"
)

// Keeper of the mint store
type Keeper struct {
	Keeper     mintkeeper.Keeper
	paramSpace paramtypes.Subspace
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	sk minttypes.StakingKeeper, ak minttypes.AccountKeeper, bk minttypes.BankKeeper,
	feeCollectorName string,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		Keeper:     mintkeeper.NewKeeper(cdc, key, paramSpace, sk, ak, bk, feeCollectorName),
		paramSpace: paramSpace,
	}
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params minttypes.Params) {
	k.paramSpace.SetParamSet(ctx, &types.Params{Params: params})
}
