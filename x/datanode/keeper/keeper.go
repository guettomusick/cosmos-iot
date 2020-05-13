package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/qonico/cosmos-iot/x/datanode/types"
)

// DataNodeKeeper - keeper of the datanode store
type DataNodeKeeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

// NewKeeper - creates a datanode keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) DataNodeKeeper {
	keeper := DataNodeKeeper{
		storeKey: key,
		cdc:      cdc,
	}
	return keeper
}

// DataNode keeper methods

// GetDataNode - gets the entire datanode metadata struct for an address
func (k DataNodeKeeper) GetDataNode(ctx sdk.Context, address sdk.AccAddress) (*types.DataNode, error) {
	store := ctx.KVStore(k.storeKey)
	if !k.IsDataNodePresent(ctx, address) {
		return nil, types.ErrInvalidDataNode
	}
	bz := store.Get([]byte(address))
	var dataNode types.DataNode
	k.cdc.MustUnmarshalBinaryBare(bz, &dataNode)
	return &dataNode, nil
}

// SetDataNode - sets the entire datanode metadata struct for an address
func (k DataNodeKeeper) SetDataNode(ctx sdk.Context, address sdk.AccAddress, dataNode *types.DataNode) {
	if dataNode.Owner.Empty() {
		return
	}

	if dataNode.ID.Empty() {
		dataNode.ID = address
	}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(address), k.cdc.MustMarshalBinaryBare(dataNode))
}

// DeleteDataNode - Deletes the entire metadata struct for an address and all related datarecords
func (k DataNodeKeeper) DeleteDataNode(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	dataNode, err := k.GetDataNode(ctx, address)
	if err != nil {
		return
	}

	for _, hash := range dataNode.Records {
		store.Delete(hash[:])
	}
	store.Delete([]byte(address))
}

// IsDataNodePresent - check if the datanode is present in the store or not
func (k DataNodeKeeper) IsDataNodePresent(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(address))
}

// GetChannels - get the channels of the datanode
func (k DataNodeKeeper) GetChannels(ctx sdk.Context, address sdk.AccAddress) (*[]types.NodeChannel, error) {
	datanode, err := k.GetDataNode(ctx, address)
	if err != nil {
		return nil, err
	}
	return &datanode.Channels, nil
}

// GetChannel - get a channel from the datanode
func (k DataNodeKeeper) GetChannel(ctx sdk.Context, address sdk.AccAddress, id string) (*types.NodeChannel, error) {
	channels, err := k.GetChannels(ctx, address)
	if err != nil {
		return nil, err
	}
	for _, c := range *channels {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, types.ErrInvalidDataNodeChannel
}

// GetRecordHashes - get the datarecord hashes belonging to the datanode
func (k DataNodeKeeper) GetRecordHashes(ctx sdk.Context, address sdk.AccAddress) (*[]types.DataRecordHash, error) {
	datanode, err := k.GetDataNode(ctx, address)
	if err != nil {
		return nil, err
	}
	return &datanode.Records, nil
}

// AddRecordHash - add a recordhash to the datanode
func (k DataNodeKeeper) AddRecordHash(ctx sdk.Context, address sdk.AccAddress, hash types.DataRecordHash) error {
	datanode, err := k.GetDataNode(ctx, address)
	if err != nil {
		return err
	}
	datanode.Records = append(datanode.Records, hash)
	k.SetDataNode(ctx, address, datanode)
	return nil
}

// AddChannel - add a new channel to the datanode
func (k DataNodeKeeper) AddChannel(ctx sdk.Context, address sdk.AccAddress, channel types.NodeChannel) error {
	datanode, err := k.GetDataNode(ctx, address)
	if err != nil {
		return err
	}
	datanode.Channels = append(datanode.Channels, channel)
	k.SetDataNode(ctx, address, datanode)
	return nil
}

// ChangeChannel - change a channel on the datanode
func (k DataNodeKeeper) ChangeChannel(ctx sdk.Context, address sdk.AccAddress, channel types.NodeChannel) error {
	datanode, err := k.GetDataNode(ctx, address)
	modified := false
	if err != nil {
		return err
	}
	for _, c := range datanode.Channels {
		if c.ID == channel.ID {
			c.Variable = channel.Variable
			modified = true
			break
		}
	}

	if !modified {
		return k.AddChannel(ctx, address, channel)
	}
	k.SetDataNode(ctx, address, datanode)
	return nil
}

// DeleteChannel - removes a channel from the datanode
func (k DataNodeKeeper) DeleteChannel(ctx sdk.Context, address sdk.AccAddress, channelID string) error {
	datanode, err := k.GetDataNode(ctx, address)
	if err != nil {
		return err
	}
	for i, c := range datanode.Channels {
		if c.ID == channelID {
			datanode.Channels[i] = datanode.Channels[len(datanode.Channels)-1]
			datanode.Channels = datanode.Channels[:len(datanode.Channels)-1]
			break
		}
	}
	k.SetDataNode(ctx, address, datanode)
	return nil
}

// SetDataNodeOwner - change the owner of the datanode
func (k DataNodeKeeper) SetDataNodeOwner(ctx sdk.Context, address sdk.AccAddress, owner sdk.AccAddress) {
	dataNode, err := k.GetDataNode(ctx, address)
	if err != nil {
		newDataNode := types.NewDataNode(address, owner)
		dataNode = &newDataNode
	} else {
		dataNode.Owner = owner
	}
	k.SetDataNode(ctx, address, dataNode)
}

// DataRecord methods

// GetDataRecord - gets the datarecord from the KVStore
func (k DataNodeKeeper) GetDataRecord(ctx sdk.Context, hash types.DataRecordHash) (*types.DataRecord, error) {
	store := ctx.KVStore(k.storeKey)
	if !k.IsDataRecordPresent(ctx, hash) {
		return nil, types.ErrInvalidDataRecord
	}
	bz := store.Get(hash[:])
	var dataRecord types.DataRecord
	k.cdc.MustUnmarshalBinaryBare(bz, &dataRecord)
	return &dataRecord, nil
}

// SetDataRecord - sets the datarecord on the KVStore
func (k DataNodeKeeper) SetDataRecord(ctx sdk.Context, dataRecord *types.DataRecord) {
	if dataRecord.DataNode.Empty() || len(dataRecord.NodeChannel.ID) == 0 {
		return
	}

	hash := types.GetDataRecordHash(dataRecord.DataNode, &dataRecord.NodeChannel, dataRecord.TimeFrame)

	store := ctx.KVStore(k.storeKey)
	store.Set(hash[:], k.cdc.MustMarshalBinaryBare(dataRecord))
}

// IsDataRecordPresent - check if the datarecord is present in the store or not
func (k DataNodeKeeper) IsDataRecordPresent(ctx sdk.Context, hash types.DataRecordHash) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(hash[:])
}

// GetLastRecords - get the latest time frame records
func (k DataNodeKeeper) GetLastRecords(ctx sdk.Context, address sdk.AccAddress, channelID string) (*[]types.Record, error) {
	channel, err := k.GetChannel(ctx, address, channelID)
	if err != nil {
		return nil, err
	}

	hash := types.GetActualDataRecordHash(address, channel)

	dataRecord, err := k.GetDataRecord(ctx, hash)
	if err != nil {
		return nil, err
	}

	return &dataRecord.Records, nil
}

// GetRecords - get records from the date time frame
func (k DataNodeKeeper) GetRecords(ctx sdk.Context, address sdk.AccAddress, channelID string, date int64) (*[]types.Record, error) {
	channel, err := k.GetChannel(ctx, address, channelID)
	if err != nil {
		return nil, err
	}

	hash := types.GetDataRecordHash(address, channel, date)

	dataRecord, err := k.GetDataRecord(ctx, hash)
	if err != nil {
		return nil, err
	}

	return &dataRecord.Records, nil
}

// AddRecord - add a new record to the time frame
func (k DataNodeKeeper) AddRecord(ctx sdk.Context, address sdk.AccAddress, channelID string, date int64, record types.Record) error {
	channel, err := k.GetChannel(ctx, address, channelID)
	if err != nil {
		return err
	}

	hash := types.GetDataRecordHash(address, channel, date)

	dataRecord, err := k.GetDataRecord(ctx, hash)
	if err != nil {
		if err == types.ErrInvalidDataRecord {
			newDataRecord := types.NewDataRecord(address, channel, date)
			dataRecord = &newDataRecord
		} else {
			return err
		}
	}

	duplicate := false
	for _, r := range dataRecord.Records {
		if r.TimeStamp == record.TimeStamp {
			duplicate = true
			break
		}
	}

	if !duplicate {
		dataRecord.Records = append(dataRecord.Records, record)
		k.SetDataRecord(ctx, dataRecord)
	}
	return nil
}

// AddRecordAtTimestamp - add a new record to the time frame using timestamp data
func (k DataNodeKeeper) AddRecordAtTimestamp(ctx sdk.Context, address sdk.AccAddress, channelID string, record types.Record) error {
	return k.AddRecord(ctx, address, channelID, int64(record.TimeStamp), record)
}

// GetIterator - get an iterator over all datanodes and datrecords
func (k DataNodeKeeper) GetIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}
