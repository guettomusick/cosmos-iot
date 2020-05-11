package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// ErrInvalidDataNode no datanode present with the given address
	ErrInvalidDataNode = sdkerrors.Register(ModuleName, 1, "no datanode present with the given address")
	// ErrInvalidDataNodeChannel no channel present with the given id on the datanode
	ErrInvalidDataNodeChannel = sdkerrors.Register(ModuleName, 2, "no channel present with the given id on the datanode")
	// ErrInvalidDataRecord no datarecord present with the given address
	ErrInvalidDataRecord = sdkerrors.Register(ModuleName, 3, "no datarecord present with the given hash")
)
