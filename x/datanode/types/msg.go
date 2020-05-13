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
	if msg.DataNode.Equals(msg.Owner) {
		return []sdk.AccAddress{msg.DataNode, msg.Owner}
	}
	return []sdk.AccAddress{msg.Owner}
}

// ChannelUpdate - channel update action definition
type ChannelUpdate struct {
	Action   string `json:"action"`   // set, delete
	ID       string `json:"id"`       // channel within the datanode
	Variable string `json:"variable"` // variable of the channel (ex. temperature, humidity)
}

// MsgUpdateChannels - changes a channel on a datanode
type MsgUpdateChannels struct {
	Owner    sdk.AccAddress  `json:"owner"`    // owner of the datanode
	DataNode sdk.AccAddress  `json:"datanode"` // datanode to update
	Updates  []ChannelUpdate `json:"channels"` // channels updates
}

// NewMsgUpdateChannels is a constructor function for MsgUpdateChannels
func NewMsgUpdateChannels(owner sdk.AccAddress, dataNode sdk.AccAddress, updates []ChannelUpdate) MsgUpdateChannels {
	return MsgUpdateChannels{
		DataNode: dataNode,
		Owner:    owner,
		Updates:  updates,
	}
}

// Route should return the name of the module
func (msg MsgUpdateChannels) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUpdateChannels) Type() string { return "set_owner" }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateChannels) ValidateBasic() error {
	if msg.DataNode.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.DataNode.String())
	}
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.DataNode.String())
	}
	if len(msg.Updates) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "no channel updates")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgUpdateChannels) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUpdateChannels) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// NewRecord - record to be added to the DataRecord time frame
type NewRecord struct {
	NodeChannelID string  `json:"channel"`   // channel within the datanode
	TimeStamp     uint32  `json:"timestamp"` // timestamp in seconds since epoch
	Value         float32 `json:"value"`     // numeric value of the record
	Misc          string  `json:"misc"`      // miscellaneous data for other non numeric records
}

// MsgAddRecords - adds new records to the datarecord time frame
type MsgAddRecords struct {
	DataNode sdk.AccAddress `json:"datanode"`
	Records  []NewRecord    `json:"records"`
}

// NewMsgAddRecords is a constructor function for MsgAddRecords
func NewMsgAddRecords(owner sdk.AccAddress, dataNode sdk.AccAddress, records []NewRecord) MsgAddRecords {
	return MsgAddRecords{
		DataNode: dataNode,
		Records:  records,
	}
}

// Route should return the name of the module
func (msg MsgAddRecords) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddRecords) Type() string { return "set_owner" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddRecords) ValidateBasic() error {
	if msg.DataNode.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.DataNode.String())
	}
	if len(msg.Records) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "no new records")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddRecords) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddRecords) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DataNode}
}
