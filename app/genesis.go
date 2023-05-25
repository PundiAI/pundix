package app

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"
	ibctypes "github.com/cosmos/ibc-go/v3/modules/core/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tendermint/tendermint/types"

	pxtypes "github.com/pundix/pundix/types"
	pundixtransfertypes "github.com/pundix/pundix/x/ibc/applications/transfer/types"
)

// GenesisState The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

func NewDefAppGenesisByDenom(stakingDenom, mintDenom string, cdc codec.JSONCodec) map[string]json.RawMessage {
	genesis := make(map[string]json.RawMessage)
	for _, b := range ModuleBasics {
		switch b.Name() {
		case stakingtypes.ModuleName:
			state := stakingtypes.DefaultGenesisState()
			state.Params.BondDenom = stakingDenom
			state.Params.MaxValidators = 20
			state.Params.UnbondingTime = time.Hour * 24 * 21
			state.Params.HistoricalEntries = 20000
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case slashingtypes.ModuleName:
			state := slashingtypes.DefaultGenesisState()
			state.Params.MinSignedPerWindow = sdk.NewDecWithPrec(5, 2)
			state.Params.SignedBlocksWindow = 20000
			state.Params.SlashFractionDoubleSign = sdk.NewDec(1).Quo(sdk.NewDec(20))
			state.Params.SlashFractionDowntime = sdk.NewDec(1).Quo(sdk.NewDec(1000))
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case distributiontypes.ModuleName:
			state := distributiontypes.DefaultGenesisState()
			state.Params.CommunityTax = sdk.NewDecWithPrec(25, 2)
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case govtypes.ModuleName:
			state := govtypes.DefaultGenesisState()
			coinOne := sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
			for i := 0; i < state.DepositParams.MinDeposit.Len(); i++ {
				state.DepositParams.MinDeposit[i].Denom = stakingDenom
				state.DepositParams.MinDeposit[i].Amount = coinOne.Mul(sdk.NewInt(10000))
			}
			state.DepositParams.MaxDepositPeriod = time.Hour * 24 * 14
			state.VotingParams.VotingPeriod = time.Hour * 24 * 14
			state.TallyParams.Quorum = sdk.NewDecWithPrec(4, 1)
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case crisistypes.ModuleName:
			state := crisistypes.DefaultGenesisState()
			coinOne := sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
			state.ConstantFee.Denom = stakingDenom
			state.ConstantFee.Amount = sdk.NewInt(13333).Mul(coinOne)
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case minttypes.ModuleName:
			state := minttypes.DefaultGenesisState()
			state.Params.MintDenom = mintDenom
			state.Params.InflationMin = sdk.NewDecWithPrec(2000, 2)
			state.Params.InflationMax = sdk.NewDecWithPrec(4000, 2)
			state.Params.GoalBonded = sdk.NewDecWithPrec(51, 2)
			state.Params.InflationRateChange = sdk.NewDecWithPrec(100, 2)
			state.Minter.Inflation = sdk.NewDecWithPrec(3000, 2)
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case banktypes.ModuleName:
			state := banktypes.DefaultGenesisState()
			state.DenomMetadata = []banktypes.Metadata{
				pxtypes.GetPURSEMetaData(mintDenom),
				pxtypes.GetPUNDIXMetaData(stakingDenom),
			}
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case ibchost.ModuleName:
			state := ibctypes.DefaultGenesisState()
			// only allowedClients tendermint
			state.ClientGenesis.Params.AllowedClients = []string{ibcexported.Tendermint}
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case ibctransfertypes.ModuleName:
			state := ibctransfertypes.DefaultGenesisState()
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case pundixtransfertypes.CompatibleModuleName:
			// ignore self-defined ibc module genesis
		default:
			genesis[b.Name()] = b.DefaultGenesis(cdc)
		}
	}
	return genesis
}

func CustomConsensusParams() *tmproto.ConsensusParams {
	result := types.DefaultConsensusParams()
	result.Block.MaxBytes = 1048576
	result.Block.MaxGas = -1
	result.Block.TimeIotaMs = 1000
	result.Evidence.MaxAgeNumBlocks = 1000000
	result.Evidence.MaxBytes = 100000
	result.Evidence.MaxAgeDuration = 172800000000000
	return result
}
