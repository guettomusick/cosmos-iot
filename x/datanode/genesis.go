package datanode

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/qonico/cosmos-iot/x/datanode/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k DataNodeKeeper, data GenesisState) {
	for _, dn := range data.DataNodes {
		k.SetDataNode(ctx, dn.ID, &dn)
	}

	for _, dr := range data.DataRecords {
		k.SetDataRecord(ctx, &dr)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k DataNodeKeeper) (data GenesisState) {
	var dataNodes []DataNode
	var dataRecords []DataRecord

	iterator := k.GetIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		address, err := sdk.AccAddressFromBech32(string(iterator.Key()))
		if err == nil {
			// DataNode
			dataNode, err := k.GetDataNode(ctx, address)
			if err != nil {
				dataNodes = append(dataNodes, *dataNode)
			}
		} else {
			// DataRecord
			var hash types.DataRecordHash
			copy(hash[:], []byte(iterator.Key())[:16])
			dataRecord, err := k.GetDataRecord(ctx, hash)
			if err != nil {
				dataRecords = append(dataRecords, *dataRecord)
			}
		}
	}
	return NewGenesisState()
}
