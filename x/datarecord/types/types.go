package types

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

// NewDataRecord returns a new DataRecord with the DataNode and the NodeChannel and empty records set
func NewDataRecord(dataNode sdk.AccAddress, channel NodeChannel) DataRecord {
	records := []Record{}
	return DataRecord{
		DataNode:    dataNode,
		NodeChannel: channel,
		Records:     records,
	}
}

// DataRecordHash returns the hash key to be used for KVStore
func DataRecordHash(DataNode sdk.AccAddress, channel NodeChannel) [16]byte {
	now := time.Now()
	// Use days since epoch as daily time frame to group records
	key := fmt.Sprintf("%s%s%s%d", string(DataNode), channel.ID, channel.Variable, now.Unix()/(24*3600))

	return md5.Sum([]byte(key))
}

// implement fmt.Stringer
func (d DataRecord) String() string {
	return strings.TrimSpace(fmt.Sprintf(`
		ID: %s
		Owner: %s
		Name: %s
	`, d.ID, d.Owner, d.Name))
}
