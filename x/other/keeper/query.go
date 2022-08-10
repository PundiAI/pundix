package keeper

import (
	// this line is used by starport scaffolding # 1
	"github.com/pundix/pundix/x/other/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGasPrice:
			var gasPrices sdk.Coins
			for _, coin := range ctx.MinGasPrices() {
				gasPrices = append(gasPrices, sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()))
			}
			return codec.MarshalJSONIndent(legacyQuerierCdc, gasPrices)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}
