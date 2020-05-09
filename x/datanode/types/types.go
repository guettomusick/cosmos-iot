package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NodeChannel holds information about the data channel of the DataNode
type NodeChannel struct {
	ID       string `json:"id,omitempty"` // id of the channel
	Variable string `json:"variable"`     // variable of the channel (ex. temperature, humidity)
}

// DataNode holds the configuration and the owner of the DataNode Device
type DataNode struct {
	ID       sdk.AccAddress `json:"id,omitempty"` // id of the datanode
	Owner    sdk.AccAddress `json:"owner"`        // account address that owns the NFT
	Name     string         `json:"name"`         // name of the datanode
	Channels []NodeChannel  `json:"channels"`     // channel definition
}

// NewDataNode returns a new DataNode with the ID
func NewDataNode(address sdk.AccAddress) DataNode {
	return DataNode{
		ID:   address,
		Name: string(address),
	}
}

// implement fmt.Stringer
func (d DataNode) String() string {
	return strings.TrimSpace(fmt.Sprintf(`
		ID: %s
		Owner: %s
		Name: %s
	`, d.ID, d.Owner, d.Name))
}
