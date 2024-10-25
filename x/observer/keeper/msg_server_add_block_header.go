package keeper

import (
	"context"
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// AddBlockHeader handles adding a block header to the store, through majority voting of observers
func (k msgServer) AddBlockHeader(goCtx context.Context, msg *types.MsgAddBlockHeader) (*types.MsgAddBlockHeaderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check authorization for this chain
	chain := k.GetSupportedChainFromChainID(ctx, msg.ChainId)
	if chain == nil {
		return nil, cosmoserrors.Wrapf(types.ErrSupportedChains, "chain id: %d", msg.ChainId)
	}

	if ok := k.IsAuthorized(ctx, msg.Creator); !ok {
		return nil, types.ErrNotAuthorizedPolicy
	}

	crosschainFlags, found := k.GetCrosschainFlags(ctx)
	if !found {
		return nil, fmt.Errorf("crosschain flags not found")
	}
	if crosschainFlags.BlockHeaderVerificationFlags == nil {
		return nil, fmt.Errorf("block header verification flags not found")
	}
	if common.IsBitcoinChain(msg.ChainId) && !crosschainFlags.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled {
		return nil, cosmoserrors.Wrapf(types.ErrBlockHeaderVerificationDisabled, "proof verification not enabled for bitcoin ,chain id: %d", msg.ChainId)
	}
	if common.IsEVMChain(msg.ChainId) && !crosschainFlags.BlockHeaderVerificationFlags.IsEthTypeChainEnabled {
		return nil, cosmoserrors.Wrapf(types.ErrBlockHeaderVerificationDisabled, "proof verification not enabled for evm ,chain id: %d", msg.ChainId)
	}

	_, found = k.GetBlockHeader(ctx, msg.BlockHash)
	if found {
		return nil, cosmoserrors.Wrap(types.ErrBlockAlreadyExist, fmt.Sprintf("block hash: %x", msg.BlockHash))
	}

	// if BlockHeaderState exists and parent block header is not found, this tx is rejected
	// if no BlockHeaderState is found, allow this vote to pass through to create and initialize
	// the Earliest/Latest height with this block header (after voting, not here)
	// if BlockHeaderState is found, check if the block height is valid
	// validate block height as it's not part of the header itself
	bhs, found := k.Keeper.GetBlockHeaderState(ctx, msg.ChainId)
	if found && bhs.EarliestHeight > 0 && bhs.EarliestHeight < msg.Height {
		pHash, err := msg.Header.ParentHash()
		if err != nil {
			return nil, cosmoserrors.Wrap(types.ErrNoParentHash, err.Error())
		}
		_, found = k.GetBlockHeader(ctx, pHash)
		if !found {
			return nil, cosmoserrors.Wrap(types.ErrNoParentHash, "parent block header not found")
		}
		if msg.Height != bhs.LatestHeight+1 {
			return nil, cosmoserrors.Wrap(types.ErrNoParentHash, fmt.Sprintf("invalid block height: wanted %d, got %d", bhs.LatestHeight+1, msg.Height))
		}
	}

	// Check timestamp
	err := msg.Header.ValidateTimestamp(ctx.BlockTime())
	if err != nil {
		return nil, cosmoserrors.Wrap(types.ErrInvalidTimestamp, err.Error())
	}

	// NOTE: error is checked in BasicValidation in msg; check again for extra caution
	pHash, err := msg.Header.ParentHash()
	if err != nil {
		return nil, cosmoserrors.Wrap(types.ErrNoParentHash, err.Error())
	}

	// add vote to ballot
	ballot, _, err := k.FindBallot(ctx, msg.Digest(), chain, types.ObservationType_InBoundTx)
	if err != nil {
		return nil, cosmoserrors.Wrap(err, "failed to find ballot")
	}
	ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, types.VoteType_SuccessObservation)
	if err != nil {
		return nil, cosmoserrors.Wrap(err, "failed to add vote to ballot")
	}
	_, isFinalized := k.CheckIfFinalizingVote(ctx, ballot)
	if !isFinalized {
		return &types.MsgAddBlockHeaderResponse{}, nil
	}

	/**
	 * Vote finalized, add block header to store
	 */
	bhs, found = k.Keeper.GetBlockHeaderState(ctx, msg.ChainId)
	if !found {
		bhs = types.BlockHeaderState{
			ChainId:         msg.ChainId,
			LatestHeight:    msg.Height,
			EarliestHeight:  msg.Height,
			LatestBlockHash: msg.BlockHash,
		}
	} else {
		if msg.Height > bhs.LatestHeight {
			bhs.LatestHeight = msg.Height
			bhs.LatestBlockHash = msg.BlockHash
		}
		if bhs.EarliestHeight == 0 {
			bhs.EarliestHeight = msg.Height
		}
	}
	k.Keeper.SetBlockHeaderState(ctx, bhs)

	bh := common.BlockHeader{
		Header:     msg.Header,
		Height:     msg.Height,
		Hash:       msg.BlockHash,
		ParentHash: pHash,
		ChainId:    msg.ChainId,
	}
	k.SetBlockHeader(ctx, bh)

	return &types.MsgAddBlockHeaderResponse{}, nil
}
