package keeper

import (
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/qonico/cosmos-iot/x/datanode/types"
)

// NewQuerier creates a new querier for datanode clients.
func NewQuerier(k DataNodeKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryDataNode:
			return queryDataNode(ctx, path[1:], req, k)
		case types.QueryRecords:
			return queryRecords(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown datanode query endpoint")
		}
	}
}

func queryDataNode(ctx sdk.Context, path []string, req abci.RequestQuery, k DataNodeKeeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	datanode, err := k.GetDataNode(ctx, address)
	if err != nil {
		return nil, err
	}

	res, err := codec.MarshalJSONIndent(k.cdc, datanode)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryRecords(ctx sdk.Context, path []string, req abci.RequestQuery, k DataNodeKeeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	date, err := strconv.ParseInt(path[2], 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	records, err := k.GetRecords(ctx, address, path[1], date)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidDataRecord, err.Error())
	}

	var resRecords (types.QueryResRecordsList)

	for _, re := range *records {
		resRecords = append(resRecords, types.QueryResRecords{
			TimeStamp: re.TimeStamp,
			Value:     re.Value,
			Misc:      re.Misc,
		})
	}
	res, err := codec.MarshalJSONIndent(k.cdc, resRecords)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
