package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/version"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/ibc-go/modules/core/02-client/types"
	"github.com/spf13/cobra"
)

func ClientUpdateProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "client-update-proposal [title] [desc] [depositAmount] [subjectClientID] [substituteClientID]",
		Args: cobra.ExactArgs(5),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s transact client-update-proposal "proposal update ibc client - 1" "proposal ..." 100000000FX 07-tendermint-0 07-tendermint-1
`, version.AppName,
		)),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sender := clientCtx.GetFromAddress()

			proposalTitle := args[0]
			proposalDesc := args[1]
			proposalDepositAmount := args[2]
			subjectClientID := args[3]
			substituteClientID := args[4]
			deposit, err := sdk.ParseCoinsNormalized(proposalDepositAmount)
			if err != nil {
				return err
			}
			proposal := types.NewClientUpdateProposal(proposalTitle, proposalDesc, subjectClientID, substituteClientID)
			msgSubmitProposal, err := govtypes.NewMsgSubmitProposal(proposal, deposit, sender)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgSubmitProposal)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
