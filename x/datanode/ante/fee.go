package ante

import (
	"fmt"

	"github.com/qonico/cosmos-iot/x/datanode/keeper"
	"github.com/qonico/cosmos-iot/x/datanode/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authAnte "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ DelegatedFeeTx = (*authTypes.StdTx)(nil) // assert StdTx implements FeeTx
)

// DelegatedFeeTx defines the interface to be implemented by Tx to use the FeeDecorators
type DelegatedFeeTx interface {
	sdk.Tx
	GetGas() uint64
	GetFee() sdk.Coins
	GetSigners() []sdk.AccAddress
}

// DelegatedDeductFeeDecorator deducts fees from the delegated account or the first signer of the tx
// If the fee payer does not have the funds to pay for the fees, return with InsufficientFunds error
// Call next AnteHandler if fees successfully deducted
// CONTRACT: Tx must implement FeeTx interface to use DelegatedDeductFeeDecorator
type DelegatedDeductFeeDecorator struct {
	ak             authKeeper.AccountKeeper
	supplyKeeper   authTypes.SupplyKeeper
	dataNodeKeeper keeper.DataNodeKeeper
}

func NewDelegatedDeductFeeDecorator(ak authKeeper.AccountKeeper, sk authTypes.SupplyKeeper, dk keeper.DataNodeKeeper) DelegatedDeductFeeDecorator {
	return DelegatedDeductFeeDecorator{
		ak:             ak,
		supplyKeeper:   sk,
		dataNodeKeeper: dk,
	}
}

func (dfd DelegatedDeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(DelegatedFeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a DelegatedFeeTx")
	}

	if addr := dfd.supplyKeeper.GetModuleAddress(authTypes.FeeCollectorName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", authTypes.FeeCollectorName))
	}

	signerAddrs := feeTx.GetSigners()
	var dataNode (*types.DataNode)

	for _, sa := range signerAddrs {
		dn, err := dfd.dataNodeKeeper.GetDataNode(ctx, sa)
		if err != nil {
			dataNode = dn
			break
		}
	}

	// Check if some DataNode signed the Transaction, if not use default DeductFeeDecorator
	if dataNode == nil {
		return authAnte.NewDeductFeeDecorator(dfd.ak, dfd.supplyKeeper).AnteHandle(ctx, tx, simulate, next)
	}

	feePayer := dataNode.Owner
	feePayerAcc := dfd.ak.GetAccount(ctx, feePayer)

	if feePayerAcc == nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	// deduct the fees
	if !feeTx.GetFee().IsZero() {
		err = authAnte.DeductFees(dfd.supplyKeeper, ctx, feePayerAcc, feeTx.GetFee())
		if err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}
