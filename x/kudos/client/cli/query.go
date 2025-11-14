package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/pavlenkotm/cosmos-kudos-module/x/kudos/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryBalance(),
		CmdQueryLeaderboard(),
	)

	return cmd
}

// CmdQueryBalance returns a CLI command handler for querying kudos balance
func CmdQueryBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance [address]",
		Short: "Query kudos balance for an address",
		Long: `Query the kudos balance for a specific address.

Example:
  kudos balance cosmos1...
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryKudosBalanceRequest{
				Address: args[0],
			}

			res, err := queryClient.KudosBalance(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryLeaderboard returns a CLI command handler for querying the kudos leaderboard
func CmdQueryLeaderboard() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leaderboard [limit]",
		Short: "Query kudos leaderboard",
		Long: `Query the kudos leaderboard showing top receivers.

Example:
  kudos leaderboard 10
`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			limit := uint32(10) // default limit
			if len(args) > 0 {
				parsedLimit, err := strconv.ParseUint(args[0], 10, 32)
				if err != nil {
					return fmt.Errorf("invalid limit: %w", err)
				}
				limit = uint32(parsedLimit)
			}

			params := &types.QueryKudosLeaderboardRequest{
				Limit: limit,
			}

			res, err := queryClient.KudosLeaderboard(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
