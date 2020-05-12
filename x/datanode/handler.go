package datanode

import (
	"fmt"

	"github.com/qonico/cosmos-iot/x/datanode/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the datanode type messages
func NewHandler(k DataNodeKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgSetOwner:
			return handleMsgSetOwner(ctx, k, msg)
		case types.MsgUpdateChannels:
			return handleMsgUpdateChannels(ctx, k, msg)
		case types.MsgAddRecords:
			return handleMsgAddRecords(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgSetOwner - handle a messsage to change owner or create a new datanode
func handleMsgSetOwner(ctx sdk.Context, k DataNodeKeeper, msg types.MsgSetOwner) (*sdk.Result, error) {
	dataNode, err := k.GetDataNode(ctx, msg.DataNode)
	if err != nil {
		// If theres no datanode, then owner must be the datanode itself
		if !msg.Owner.Equals(msg.DataNode) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner - owner must be the same as datanode for datanode creation")
		}
	} else if !dataNode.Owner.Equals(msg.Owner) {
		// only owner can reassign owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner - existing datanode and owner don't match")
	}

	k.SetDataNodeOwner(ctx, msg.DataNode, msg.NewOwner)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgUpdateChannels - handle a messsage to update channels definition
func handleMsgUpdateChannels(ctx sdk.Context, k DataNodeKeeper, msg types.MsgUpdateChannels) (*sdk.Result, error) {
	dataNode, err := k.GetDataNode(ctx, msg.DataNode)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "Incorrect DataNode - not defined")
	}
	if !dataNode.Owner.Equals(msg.Owner) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner - existing datanode and owner don't match")
	}

	for _, ch := range msg.Updates {
		switch ch.Action {
		case "set":
			channel := types.NodeChannel{
				ID:       ch.ID,
				Variable: ch.Variable,
			}
			k.ChangeChannel(ctx, msg.DataNode, channel)
			break
		case "delete":
			k.DeleteChannel(ctx, msg.DataNode, ch.ID)
			break
		}
	}
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgAddRecords - handle a messsage to add records to persist
func handleMsgAddRecords(ctx sdk.Context, k DataNodeKeeper, msg types.MsgAddRecords) (*sdk.Result, error) {
	_, err := k.GetDataNode(ctx, msg.DataNode)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "Incorrect DataNode - not defined")
	}

	for _, re := range msg.Records {
		record := types.Record{
			TimeStamp: re.TimeStamp,
			Value:     re.Value,
			Misc:      re.Misc,
		}
		k.AddRecordAtTimestamp(ctx, msg.DataNode, re.NodeChannelID, record)
	}
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
