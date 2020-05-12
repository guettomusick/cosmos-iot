package types

import (
	"encoding/json"
)

// Query endpoints supported by the datanode querier
const (
	QueryDataNode = "datanode"
	QueryRecords  = "records"
)

// QueryResRecords - queries result payload for a single record
type QueryResRecords struct {
	TimeStamp uint32  `json:"ts"`    // timestamp in seconds since epoch
	Value     float32 `json:"value"` // numeric value of the record
	Misc      string  `json:"misc"`  // miscellaneous data for other non numeric records
}

// QueryResRecordsList - queries result payload for datarecords within time frame
type QueryResRecordsList []QueryResRecords

// implement fmt.Stringer
func (r QueryResRecordsList) String() string {
	res, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(res)
}
