package types

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DataRecordHash is the hash key of the records time frame
type DataRecordHash [16]byte

// NodeChannel holds information about the data channel of the DataNode
type NodeChannel struct {
	ID       string `json:"id,omitempty"` // id of the channel
	Variable string `json:"variable"`     // variable of the channel (ex. temperature, humidity)
}

// DataNode holds the configuration and the owner of the DataNode Device
type DataNode struct {
	ID       sdk.AccAddress   `json:"id,omitempty"` // id of the datanode
	Owner    sdk.AccAddress   `json:"owner"`        // account address that owns the DataNode
	Name     string           `json:"name"`         // name of the datanode
	Channels []NodeChannel    `json:"channels"`     // channel definition
	Records  []DataRecordHash `json:"records"`      // datarecords associated to this DataNode
}

// Record holds a single record from the DataNode device
type Record struct {
	TimeStamp uint32  `json:"t"` // timestamp in seconds since epoch
	Value     float32 `json:"v"` // numeric value of the record
	Misc      string  `json:"m"` // miscellaneous data for other non numeric records
}

// DataRecord is a time frame package of records
type DataRecord struct {
	DataNode    sdk.AccAddress `json:"datanode"` // datanode which push the records
	NodeChannel NodeChannel    `json:"channel"`  // channel within the datanode
	Records     []Record       `json:"records"`  // records of the timerange
}

// NewDataNode returns a new DataNode with the ID
func NewDataNode(address sdk.AccAddress, owner sdk.AccAddress) DataNode {
	return DataNode{
		ID:    address,
		Owner: owner,
		Name:  string(address),
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

// NewDataRecord returns a new DataRecord with the DataNode and the NodeChannel and empty records set
func NewDataRecord(dataNode sdk.AccAddress, channel *NodeChannel) DataRecord {
	records := []Record{}
	return DataRecord{
		DataNode:    dataNode,
		NodeChannel: *channel,
		Records:     records,
	}
}

// GetActualDataRecordHash returns the hash key to be used for KVStore at actual time
func GetActualDataRecordHash(dataNode sdk.AccAddress, channel *NodeChannel) DataRecordHash {
	now := time.Now()
	return GetDataRecordHash(dataNode, channel, now.Unix())
}

// GetDataRecordHash returns the hash key to be used for KVStore
func GetDataRecordHash(dataNode sdk.AccAddress, channel *NodeChannel, date int64) DataRecordHash {
	// Use days since epoch as daily time frame to group records
	key := fmt.Sprintf("%s%s%s%d", string(dataNode), channel.ID, channel.Variable, date/(24*3600))

	return md5.Sum([]byte(key))
}

// implement fmt.Stringer
func (r DataRecord) String() string {
	return strings.TrimSpace(fmt.Sprintf(`
		DataNode: %s
		Channel: %s:%s
		Records: %d
		From: %d
		To: %d
	`, string(r.DataNode), r.NodeChannel.ID, r.NodeChannel.Variable, len(r.Records), r.Records[0].TimeStamp, r.Records[len(r.Records)-1].TimeStamp))
}
