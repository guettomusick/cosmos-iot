package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetOwner{}, "datanode/SetOwner", nil)
	cdc.RegisterConcrete(MsgUpdateChannels{}, "datanode/UpdateChannels", nil)
	cdc.RegisterConcrete(MsgAddRecords{}, "datanode/AddRecords", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
