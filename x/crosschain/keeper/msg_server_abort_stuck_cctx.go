package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

const (
	// AbortMessage is the message to abort a stuck CCTX
	AbortMessage = "CCTX aborted with admin cmd"
)

// AbortStuckCCTX aborts a stuck CCTX
// Authorized: admin policy group 2
func (k msgServer) AbortStuckCCTX(
	goCtx context.Context,
	msg *types.MsgAbortStuckCCTX,
) (*types.MsgAbortStuckCCTXResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if authorized
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupAdmin) {
		return nil, observertypes.ErrNotAuthorized
	}

	// check if the cctx exists
	cctx, found := k.GetCrossChainTx(ctx, msg.CctxIndex)
	if !found {
		return nil, types.ErrCannotFindCctx
	}

	// check if the cctx is pending
	isPending := cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound ||
		cctx.CctxStatus.Status == types.CctxStatus_PendingInbound ||
		cctx.CctxStatus.Status == types.CctxStatus_PendingRevert
	if !isPending {
		return nil, types.ErrStatusNotPending
	}

	cctx.CctxStatus = &types.Status{
		Status:        types.CctxStatus_Aborted,
		StatusMessage: AbortMessage,
	}

	k.SetCrossChainTx(ctx, cctx)

	return &types.MsgAbortStuckCCTXResponse{}, nil
}
