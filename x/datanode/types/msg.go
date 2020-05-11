package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgSetOwner change the owner of a DataNode or creates a new one if doesn't exist
type MsgSetOwner struct {
	DataNode sdk.AccAddress `json:"datanode"`
	Owner    sdk.AccAddress `json:"owner"`
	NewOwner sdk.AccAddress `json:"newowner"`
	Name     string         `json:"name"`
}

// NewMsgSetOwner is a constructor function for MsgSetOwner
func NewMsgSetOwner(dataNode sdk.AccAddress, owner sdk.AccAddress, name string) MsgSetOwner {
	return MsgSetOwner{
		DataNode: dataNode,
		Owner:    owner,
		Name:     name,
	}
}

// Route should return the name of the module
func (msg MsgSetOwner) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetOwner) Type() string { return "set_owner" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetOwner) ValidateBasic() error {
	if msg.DataNode.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.DataNode.String())
	}
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.DataNode.String())
	}
	if msg.NewOwner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.NewOwner.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetOwner) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetOwner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
