package datanode

import (
	"github.com/qonico/cosmos-iot/x/datanode/ante"
	"github.com/qonico/cosmos-iot/x/datanode/keeper"
	"github.com/qonico/cosmos-iot/x/datanode/types"
)

const (
	// TODO: define constants that you would like exposed from your module

	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	QuerierRoute      = types.QuerierRoute
)

var (
	// functions aliases
	NewKeeper                        = keeper.NewKeeper
	NewQuerier                       = keeper.NewQuerier
	RegisterCodec                    = types.RegisterCodec
	NewGenesisState                  = types.NewGenesisState
	DefaultGenesisState              = types.DefaultGenesisState
	ValidateGenesis                  = types.ValidateGenesis
	NewDelegatedDeductFeeAnteHandler = ante.NewDelegatedDeductFeeAnteHandler

	// variable aliases
	ModuleCdc = types.ModuleCdc
	// TODO: Fill out variable aliases
)

type (
	DataNodeKeeper = keeper.DataNodeKeeper
	GenesisState   = types.GenesisState
	Params         = types.Params
	DataNode       = types.DataNode
	DataRecord     = types.DataRecord
)
