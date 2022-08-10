package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pundix/pundix/x/ibc/applications/transfer/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
)

// GetCmdQueryDenomTrace defines the command to query a a denomination trace from a given hash.
func GetCmdQueryDenomTrace() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denom-trace [hash]",
		Short:   "Query the denom trace info from a given trace hash",
		Long:    "Query the denom trace info from a given trace hash",
		Example: fmt.Sprintf("%s query ibc-transfer denom-trace [hash]", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryDenomTraceRequest{
				Hash: args[0],
			}

			res, err := queryClient.DenomTrace(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryDenomTraces defines the command to query all the denomination trace infos
// that this chain mantains.
func GetCmdQueryDenomTraces() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denom-traces",
		Short:   "Query the trace info for all token denominations",
		Long:    "Query the trace info for all token denominations",
		Example: fmt.Sprintf("%s query ibc-transfer denom-traces", version.AppName),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryDenomTracesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.DenomTraces(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "denominations trace")

	return cmd
}

// GetCmdParams returns the command handler for ibc-transfer parameter querying.
func GetCmdParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query the current ibc-transfer parameters",
		Long:    "Query the current ibc-transfer parameters",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("%s query ibc-transfer params", version.AppName),
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, _ := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			return clientCtx.PrintProto(res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdDenomToIBcDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denom-convert",
		Short:   "Covert denom to ibc denom",
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf("%s query ibc-transfer denom-convert transfer/channel-0/0x2170ed0880ac9a755fd29b2688956bd959f933f8", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			denomTrace := types.ParseDenomTrace(args[0])

			type output struct {
				Prefix   string
				Denom    string
				IBCDenom string
			}

			marshal, err := json.Marshal(output{
				Prefix:   denomTrace.GetPrefix(),
				Denom:    denomTrace.GetBaseDenom(),
				IBCDenom: denomTrace.IBCDenom(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintBytes(marshal)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdGetEscrowAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "escrow-address",
		Short:   "get escrow address",
		Args:    cobra.ExactArgs(2),
		Example: fmt.Sprintf("%s query ibc-transfer escrow-address transfer channel-0", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			portId := args[0]
			channelId := args[1]
			escrowAddress := types.GetEscrowAddress(portId, channelId)
			return clientCtx.PrintObjectLegacy(map[string]interface{}{
				"port_id":        portId,
				"channel_id":     channelId,
				"escrow_address": escrowAddress.String(),
			})
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
