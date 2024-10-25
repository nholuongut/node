package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// SetCctxAndNonceToCctxAndInTxHashToCctx does the following things in one function:
// 1. set the cctx in the store
// 2. set the mapping inTxHash -> cctxIndex , one inTxHash can be connected to multiple cctxindex
// 3. set the mapping nonce => cctx
// 4. update the zeta accounting
func (k Keeper) SetCctxAndNonceToCctxAndInTxHashToCctx(ctx sdk.Context, cctx types.CrossChainTx) {
	k.SetCrossChainTx(ctx, cctx)

	// set mapping inTxHash -> cctxIndex
	in, _ := k.GetInTxHashToCctx(ctx, cctx.InboundTxParams.InboundTxObservedHash)
	in.InTxHash = cctx.InboundTxParams.InboundTxObservedHash
	found := false
	for _, cctxIndex := range in.CctxIndex {
		if cctxIndex == cctx.Index {
			found = true
			break
		}
	}
	if !found {
		in.CctxIndex = append(in.CctxIndex, cctx.Index)
	}
	k.SetInTxHashToCctx(ctx, in)

	tss, found := k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return
	}
	// set mapping nonce => cctxIndex
	if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound || cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
		k.GetObserverKeeper().SetNonceToCctx(ctx, observerTypes.NonceToCctx{
			ChainId: cctx.GetCurrentOutTxParam().ReceiverChainId,
			// #nosec G701 always in range
			Nonce:     int64(cctx.GetCurrentOutTxParam().OutboundTxTssNonce),
			CctxIndex: cctx.Index,
			Tss:       tss.TssPubkey,
		})
	}
	if cctx.CctxStatus.Status == types.CctxStatus_Aborted && cctx.GetCurrentOutTxParam().CoinType == common.CoinType_Zeta {
		k.AddZetaAbortedAmount(ctx, GetAbortedAmount(cctx))
	}
}

// SetCrossChainTx set a specific send in the store from its index
func (k Keeper) SetCrossChainTx(ctx sdk.Context, cctx types.CrossChainTx) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.SendKey))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	b := k.cdc.MustMarshal(&cctx)
	store.Set(types.KeyPrefix(cctx.Index), b)
}

// GetCrossChainTx returns a send from its index
func (k Keeper) GetCrossChainTx(ctx sdk.Context, index string) (val types.CrossChainTx, found bool) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.SendKey))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllCrossChainTx(ctx sdk.Context) (list []types.CrossChainTx) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.SendKey))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)

	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CrossChainTx
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

// RemoveCrossChainTx removes a send from the store
func (k Keeper) RemoveCrossChainTx(ctx sdk.Context, index string) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.SendKey))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	store.Delete(types.KeyPrefix(index))
}

func (k Keeper) CreateNewCCTX(
	ctx sdk.Context,
	msg *types.MsgVoteOnObservedInboundTx,
	index string,
	tssPubkey string,
	s types.CctxStatus,
	senderChainID,
	receiverChainID int64,
) types.CrossChainTx {
	if msg.TxOrigin == "" {
		msg.TxOrigin = msg.Sender
	}
	inboundParams := &types.InboundTxParams{
		Sender:                          msg.Sender,
		SenderChainId:                   senderChainID,
		TxOrigin:                        msg.TxOrigin,
		Asset:                           msg.Asset,
		Amount:                          msg.Amount,
		CoinType:                        msg.CoinType,
		InboundTxObservedHash:           msg.InTxHash,
		InboundTxObservedExternalHeight: msg.InBlockHeight,
		InboundTxFinalizedZetaHeight:    0,
		InboundTxBallotIndex:            index,
	}

	outBoundParams := &types.OutboundTxParams{
		Receiver:                         msg.Receiver,
		ReceiverChainId:                  receiverChainID,
		OutboundTxHash:                   "",
		OutboundTxTssNonce:               0,
		OutboundTxGasLimit:               msg.GasLimit,
		OutboundTxGasPrice:               "",
		OutboundTxBallotIndex:            "",
		OutboundTxObservedExternalHeight: 0,
		CoinType:                         msg.CoinType, // FIXME: is this correct?
		Amount:                           sdk.NewUint(0),
		TssPubkey:                        tssPubkey,
	}
	status := &types.Status{
		Status:              s,
		StatusMessage:       "",
		LastUpdateTimestamp: ctx.BlockHeader().Time.Unix(),
		IsAbortRefunded:     false,
	}
	newCctx := types.CrossChainTx{
		Creator:          msg.Creator,
		Index:            index,
		ZetaFees:         math.ZeroUint(),
		RelayedMessage:   msg.Message,
		CctxStatus:       status,
		InboundTxParams:  inboundParams,
		OutboundTxParams: []*types.OutboundTxParams{outBoundParams},
	}
	return newCctx
}
