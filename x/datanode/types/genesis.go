package types

import "fmt"

// GenesisState - all datanode state that must be provided at genesis
type GenesisState struct {
	DataNodes   []DataNode   `json:"datanodes"`
	DataRecords []DataRecord `json:"datarecords"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState() GenesisState {
	return GenesisState{
		DataNodes:   nil,
		DataRecords: nil,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		DataNodes:   []DataNode{},
		DataRecords: []DataRecord{},
	}
}

// ValidateGenesis validates the datanode genesis parameters
func ValidateGenesis(data GenesisState) error {
	for _, dn := range data.DataNodes {
		if dn.ID == nil {
			return fmt.Errorf("invalid DataNode: Owner: %s. Error: Missing ID", dn.Owner)
		}
		if dn.Owner == nil {
			return fmt.Errorf("invalid DataNode: ID: %s. Error: Missing Owner", dn.ID)
		}
	}

	for _, dr := range data.DataRecords {
		if dr.DataNode == nil {
			return fmt.Errorf("invalid DataRecord: ChannelID: %s:%s. Error: Missing DataNode", dr.NodeChannel.ID, dr.NodeChannel.Variable)
		}
		if len(dr.NodeChannel.ID) == 0 {
			return fmt.Errorf("invalid DataRecord: DataNode: %s. Error: Missing ChannelID", dr.DataNode)
		}
		if len(dr.Records) == 0 {
			return fmt.Errorf("invalid DataRecord: DataNode: %s. Error: No Records", dr.DataNode)
		}
	}
	return nil
}
