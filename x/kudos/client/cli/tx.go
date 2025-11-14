package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/pavlenkotm/cosmos-kudos-module/x/kudos/types"
)

const (
	FlagComment = "comment"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdSendKudos())

	return cmd
}

// CmdSendKudos returns a CLI command handler for sending kudos
func CmdSendKudos() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [to_address] [amount]",
		Short: "Send kudos to another address",
		Long: `Send kudos to another address with an optional comment.

Example:
  kudos send cosmos1... 10 --comment "Thanks for the code review!"
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			toAddress := args[0]
			amount, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}

			comment, err := cmd.Flags().GetString(FlagComment)
			if err != nil {
				return err
			}

			// Validate comment length
			if len(comment) > 140 {
				return fmt.Errorf("comment exceeds 140 characters")
			}

			msg := &types.MsgSendKudos{
				FromAddress: clientCtx.GetFromAddress().String(),
				ToAddress:   toAddress,
				Amount:      amount,
				Comment:     comment,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagComment, "", "Comment for the kudos (max 140 characters)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
