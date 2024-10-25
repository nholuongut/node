package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// UpdateCrosschainFlags updates the crosschain related flags.
//
// Aurthorized: admin policy group 1 (except enabling/disabled
// inbounds/outbounds and gas price increase), admin policy group 2 (all).
func (k msgServer) UpdateCrosschainFlags(goCtx context.Context, msg *types.MsgUpdateCrosschainFlags) (*types.MsgUpdateCrosschainFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, msg.GetRequiredPolicyType()) {
		return &types.MsgUpdateCrosschainFlagsResponse{}, types.ErrNotAuthorizedPolicy
	}

	// check if the value exists
	flags, isFound := k.GetCrosschainFlags(ctx)
	if !isFound {
		flags = *types.DefaultCrosschainFlags()
	}

	// update values
	flags.IsInboundEnabled = msg.IsInboundEnabled
	flags.IsOutboundEnabled = msg.IsOutboundEnabled

	if msg.GasPriceIncreaseFlags != nil {
		flags.GasPriceIncreaseFlags = msg.GasPriceIncreaseFlags
	}

	if msg.BlockHeaderVerificationFlags != nil {
		flags.BlockHeaderVerificationFlags = msg.BlockHeaderVerificationFlags
	}

	k.SetCrosschainFlags(ctx, flags)

	err := ctx.EventManager().EmitTypedEvents(&types.EventCrosschainFlagsUpdated{
		MsgTypeUrl:                   sdk.MsgTypeURL(&types.MsgUpdateCrosschainFlags{}),
		IsInboundEnabled:             msg.IsInboundEnabled,
		IsOutboundEnabled:            msg.IsOutboundEnabled,
		GasPriceIncreaseFlags:        msg.GasPriceIncreaseFlags,
		BlockHeaderVerificationFlags: msg.BlockHeaderVerificationFlags,
		Signer:                       msg.Creator,
	})
	if err != nil {
		ctx.Logger().Error("Error emitting EventCrosschainFlagsUpdated :", err)
	}

	return &types.MsgUpdateCrosschainFlagsResponse{}, nil
}
