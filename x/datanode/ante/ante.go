package ante

import (
	"github.com/qonico/cosmos-iot/x/datanode/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authAnte "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NewDelegatedDeductFeeAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from deducted wallet or the first
// signer.
func NewDelegatedDeductFeeAnteHandler(ak authKeeper.AccountKeeper, supplyKeeper authTypes.SupplyKeeper, dk keeper.DataNodeKeeper, sigGasConsumer authAnte.SignatureVerificationGasConsumer) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		authAnte.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		authAnte.NewMempoolFeeDecorator(),
		authAnte.NewValidateBasicDecorator(),
		authAnte.NewValidateMemoDecorator(ak),
		authAnte.NewConsumeGasForTxSizeDecorator(ak),
		authAnte.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		authAnte.NewValidateSigCountDecorator(ak),
		NewDelegatedDeductFeeDecorator(ak, supplyKeeper, dk),
		authAnte.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		authAnte.NewSigVerificationDecorator(ak),
		authAnte.NewIncrementSequenceDecorator(ak), // innermost AnteDecorator
	)
}
