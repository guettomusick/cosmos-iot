package datanode

import (
	"fmt"

	"github.com/qonico/cosmos-iot/x/datanode/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the datanode type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgSetOwner:
			return handleMsgSetOwner(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handle a messsage to change owner or create a new datanode
func handleMsgSetOwner(ctx sdk.Context, k Keeper, msg types.MsgSetOwner) (*sdk.Result, error) {
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
