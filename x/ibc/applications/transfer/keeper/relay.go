package keeper

import (
	"fmt"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	coretypes "github.com/cosmos/ibc-go/v3/modules/core/types"

	pxtypes "github.com/pundix/pundix/types"
	"github.com/pundix/pundix/x/ibc/applications/transfer/types"
)

// SendTransfer handles transfer sending logic. There are 2 possible cases:
//
// 1. Sender chain is acting as the source zone. The coins are transferred
// to an escrow address (i.e locked) on the sender chain and then transferred
// to the receiving chain through IBC TAO logic. It is expected that the
// receiving chain will mint vouchers to the receiving address.
//
// 2. Sender chain is acting as the sink zone. The coins (vouchers) are burned
// on the sender chain and then transferred to the receiving chain though IBC
// TAO logic. It is expected that the receiving chain, which had previously
// sent the original denomination, will unescrow the fungible token and send
// it to the receiving address.
//
// Another way of thinking of source and sink zones is through the token's
// timeline. Each send to any chain other than the one it was previously
// received from is a movement forwards in the token's timeline. This causes
// trace to be added to the token's history and the destination port and
// destination channel to be prefixed to the denomination. In these instances
// the sender chain is acting as the source zone. When the token is sent back
// to the chain it previously received from, the prefix is removed. This is
// a backwards movement in the token's timeline and the sender chain
// is acting as the sink zone.
//
// Example:
// These steps of transfer occur: A -> B -> C -> A -> C -> B -> A
//
// 1. A -> B : sender chain is source zone. Denom upon receiving: 'B/denom'
// 2. B -> C : sender chain is source zone. Denom upon receiving: 'C/B/denom'
// 3. C -> A : sender chain is source zone. Denom upon receiving: 'A/C/B/denom'
// 4. A -> C : sender chain is sink zone. Denom upon receiving: 'C/B/denom'
// 5. C -> B : sender chain is sink zone. Denom upon receiving: 'B/denom'
// 6. B -> A : sender chain is sink zone. Denom upon receiving: 'denom'
func (k Keeper) SendTransfer(
	ctx sdk.Context,
	sourcePort,
	sourceChannel string,
	token sdk.Coin,
	sender sdk.AccAddress,
	receiver string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	router string,
	fee sdk.Coin,
	memo string,
) error {
	_, err := k.sendTransfer(ctx, sourcePort, sourceChannel, token, sender, receiver, timeoutHeight, timeoutTimestamp, router, fee, memo)
	return err
}

func (k Keeper) sendTransfer(ctx sdk.Context, sourcePort, sourceChannel string, token sdk.Coin, sender sdk.AccAddress,
	receiver string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, router string, fee sdk.Coin, memo string,
) (uint64, error) {
	if !k.GetSendEnabled(ctx) {
		return 0, transfertypes.ErrSendDisabled
	}

	if !k.bankKeeper.IsSendEnabledCoin(ctx, token) {
		return 0, sdkerrors.Wrapf(transfertypes.ErrSendDisabled, "%s transfers are currently disabled", token.Denom)
	}

	if k.bankKeeper.BlockedAddr(sender) {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to send funds", sender)
	}

	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return 0, sdkerrors.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return 0, sdkerrors.Wrapf(channeltypes.ErrSequenceSendNotFound, "source port: %s, source channel: %s", sourcePort, sourceChannel)
	}

	// begin createOutgoingPacket logic
	// See spec for this logic: https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#packet-relay
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	// NOTE: denomination and hex hash correctness checked during msg.ValidateBasic
	fullDenomPath := token.Denom

	var err error

	// deconstruct the token denomination into the denomination trace info
	// to determine if the sender is the source chain
	if strings.HasPrefix(token.Denom, "ibc/") {
		fullDenomPath, err = k.DenomPathFromHash(ctx, token.Denom)
		if err != nil {
			return 0, err
		}
	}

	labels := []metrics.Label{
		telemetry.NewLabel(coretypes.LabelDestinationPort, destinationPort),
		telemetry.NewLabel(coretypes.LabelDestinationChannel, destinationChannel),
	}

	packetData := types.NewFungibleTokenPacketData(
		fullDenomPath, token.Amount.String(), sender.String(), receiver, router, fee.Amount.String(),
	)

	packetData.Memo = memo
	// If the router address is specified, the number of token + fee is deducted
	if router != "" {
		token = token.Add(sdk.NewCoin(token.Denom, fee.Amount))
	}
	// NOTE: SendTransfer simply sends the denomination as it exists on its own
	// chain inside the packet data. The receiving chain will perform denom
	// prefixing as necessary.
	if transfertypes.SenderChainIsSource(sourcePort, sourceChannel, fullDenomPath) {
		labels = append(labels, telemetry.NewLabel(coretypes.LabelSource, "true"))

		// create the escrow address for the tokens
		escrowAddress := transfertypes.GetEscrowAddress(sourcePort, sourceChannel)

		// escrow source tokens. It fails if balance insufficient.
		if err = k.bankKeeper.SendCoins(ctx, sender, escrowAddress, sdk.NewCoins(token)); err != nil {
			return 0, err
		}

	} else {
		labels = append(labels, telemetry.NewLabel(coretypes.LabelSource, "false"))

		// transfer the coins to the module account and burn them
		if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, transfertypes.ModuleName, sdk.NewCoins(token)); err != nil {
			return 0, err
		}

		if err = k.bankKeeper.BurnCoins(ctx, transfertypes.ModuleName, sdk.NewCoins(token)); err != nil {
			// NOTE: should not happen as the module account was
			// retrieved on the step above and it has enough balace
			// to burn.
			panic(fmt.Sprintf("cannot burn coins after a successful send to a module account: %v", err))
		}
	}

	// hard fork: compatibility fx ibc transfer fee
	var packetDateBytes []byte
	if ctx.BlockHeight() >= pxtypes.CompatibilityIBCTransferHeight() {
		packetDateBytes = packetData.CompatibilityGetBytes()
	} else {
		packetDateBytes = packetData.GetBytes()
	}
	packet := channeltypes.NewPacket(
		packetDateBytes,
		sequence,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		timeoutHeight,
		timeoutTimestamp,
	)

	if err = k.ics4Wrapper.SendPacket(ctx, channelCap, packet); err != nil {
		return 0, err
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"ibc", transfertypes.ModuleName, "send"},
			1,
			labels,
		)
	}()

	return sequence, nil
}

// OnRecvPacket processes a cross chain fungible token transfer. If the
// sender chain is the source of minted tokens then vouchers will be minted
// and sent to the receiving address. Otherwise if the sender chain is sending
// back tokens this chain originally transferred to it, the tokens are
// unescrowed and sent to the receiving address.
func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data types.FungibleTokenPacketData) error {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return err
	}

	if len(data.Router) > 0 {
		return sdkerrors.Wrapf(types.ErrNotSupportRouter, "router:(%s)", data.Router)
	}

	if data.Fee != sdk.ZeroInt().String() {
		return sdkerrors.Wrapf(types.ErrNotSupportFee, "fee:(%s)", data.Fee)
	}

	receiver, transferAmount, feeAmount, err := parseReceiveAndAmountByPacket(data)
	if err != nil {
		return err
	}

	receiveAmount := transferAmount.Add(feeAmount)
	packetData := transfertypes.NewFungibleTokenPacketData(data.GetDenom(), receiveAmount.String(), data.GetSender(), receiver.String())
	packetData.Memo = data.Memo
	if err = k.Keeper.OnRecvPacket(ctx, packet, packetData); err != nil {
		return err
	}

	receiveDenom := parseIBCCoinDenom(packet, data.GetDenom())

	receiveCoin := sdk.NewCoin(receiveDenom, receiveAmount)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeReceive,
			sdk.NewAttribute(transfertypes.AttributeKeyReceiver, receiver.String()),
			sdk.NewAttribute(transfertypes.AttributeKeyAmount, receiveCoin.String()),
		),
	)

	return nil
}

// OnAcknowledgementPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain. If the acknowledgement
// was a success then nothing occurs. If the acknowledgement failed, then
// the sender is refunded their tokens using the refundPacketToken function.
func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, data types.FungibleTokenPacketData, ack channeltypes.Acknowledgement) error {
	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		amount, fee, err := parseAmountAndFeeByPacket(data)
		if err != nil {
			return err
		}
		ibcPacketData := data.ToIBCPacketData()
		ibcPacketData.Amount = amount.Add(fee).String()
		return k.Keeper.OnAcknowledgementPacket(ctx, packet, ibcPacketData, ack)
	default:
		// the acknowledgement succeeded on the receiving chain so nothing
		// needs to be executed and no error needs to be returned
		return nil
	}
}

// OnTimeoutPacket refunds the sender since the original packet sent was
// never received and has been timed out.
func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data types.FungibleTokenPacketData) error {
	amount, fee, err := parseAmountAndFeeByPacket(data)
	if err != nil {
		return err
	}
	ibcPacketData := data.ToIBCPacketData()
	ibcPacketData.Amount = amount.Add(fee).String()
	return k.Keeper.OnTimeoutPacket(ctx, packet, ibcPacketData)
}
