package cli

import (
	"bufio"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/qonico/cosmos-iot/x/datanode/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	datanodeTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	datanodeTxCmd.AddCommand(flags.PostCommands(
		GetCmdSetOwner(cdc),
		GetCmdUpdateChannels(cdc),
		GetCmdAddRecords(cdc),
	)...)

	return datanodeTxCmd
}

// GetCmdSetOwner is the CLI command for sending a BuyName transaction
func GetCmdSetOwner(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-owner [datanode] [owner] [newowner] [name]",
		Short: "set owner of datanode or register a new one",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			datanode, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			owner, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			newOwner, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetOwner(datanode, owner, newOwner, args[3])
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdUpdateChannels is the CLI command for sending a BuyName transaction
func GetCmdUpdateChannels(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "update-channels [owner] [datanode] [channels]",
		Short: "update channels of datanode",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			datanode, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			var channels ([]types.ChannelUpdate)
			cdc.MustUnmarshalJSON([]byte(args[2]), &channels)

			msg := types.NewMsgUpdateChannels(owner, datanode, channels)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdAddRecords is the CLI command for sending a BuyName transaction
func GetCmdAddRecords(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "add-records [datanode] [records]",
		Short: "add records to data record time frame",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			datanode, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			var records ([]types.NewRecord)
			cdc.MustUnmarshalJSON([]byte(args[1]), &records)

			msg := types.NewMsgAddRecords(datanode, records)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
