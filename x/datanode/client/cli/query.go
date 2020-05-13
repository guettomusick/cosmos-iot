package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/qonico/cosmos-iot/x/datanode/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group datanode queries under a subcommand
	datanodeQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	datanodeQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdDataNode(types.StoreKey, cdc),
			GetCmdRecords(types.StoreKey, cdc),
		)...,
	)

	return datanodeQueryCmd
}

// GetCmdDataNode queries information about a datanode
func GetCmdDataNode(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "datanode [address]",
		Short: "datanode address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/datanode/%s", queryRoute, address), nil)
			if err != nil {
				fmt.Printf("could not get datanode - %s \n", address)
				return nil
			}

			var out types.DataNode
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdRecords queries information about records on a time frame
func GetCmdRecords(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "records [address] [channelID] [date]",
		Short: "records address channelID date",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]
			channelID := args[1]
			date := args[2]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/records/%s/%s/%s", queryRoute, address, channelID, date), nil)
			if err != nil {
				fmt.Printf("could not get records on - %s %s %s \n", address, channelID, date)
				return nil
			}

			var out types.QueryResRecordsList
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
