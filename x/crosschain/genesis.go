package crosschain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// InitGenesis initializes the crosschain module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Params
	k.SetParams(ctx, genState.Params)

	k.SetZetaAccounting(ctx, genState.ZetaAccounting)
	// Set all the outTxTracker
	for _, elem := range genState.OutTxTrackerList {
		k.SetOutTxTracker(ctx, elem)
	}

	// Set all the inTxTracker
	for _, elem := range genState.InTxTrackerList {
		k.SetInTxTracker(ctx, elem)
	}

	// Set all the inTxHashToCctx
	for _, elem := range genState.InTxHashToCctxList {
		k.SetInTxHashToCctx(ctx, elem)
	}

	// Set all the gasPrice
	for _, elem := range genState.GasPriceList {
		if elem != nil {
			k.SetGasPrice(ctx, *elem)
		}
	}

	// Set all the chain nonces

	// Set all the last block heights
	for _, elem := range genState.LastBlockHeightList {
		if elem != nil {
			k.SetLastBlockHeight(ctx, *elem)
		}
	}

	// Set all the cross-chain txs
	for _, elem := range genState.CrossChainTxs {
		if elem != nil {
			k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, *elem)
		}
	}
	for _, elem := range genState.FinalizedInbounds {
		k.SetFinalizedInbound(ctx, elem)
	}
}

// ExportGenesis returns the crosschain module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var genesis types.GenesisState

	genesis.Params = k.GetParams(ctx)
	genesis.OutTxTrackerList = k.GetAllOutTxTracker(ctx)
	genesis.InTxHashToCctxList = k.GetAllInTxHashToCctx(ctx)
	genesis.InTxTrackerList = k.GetAllInTxTracker(ctx)

	// Get all gas prices
	gasPriceList := k.GetAllGasPrice(ctx)
	for _, elem := range gasPriceList {
		elem := elem
		genesis.GasPriceList = append(genesis.GasPriceList, &elem)
	}

	// Get all last block heights
	lastBlockHeightList := k.GetAllLastBlockHeight(ctx)
	for _, elem := range lastBlockHeightList {
		elem := elem
		genesis.LastBlockHeightList = append(genesis.LastBlockHeightList, &elem)
	}

	// Get all send
	sendList := k.GetAllCrossChainTx(ctx)
	for _, elem := range sendList {
		elem := elem
		genesis.CrossChainTxs = append(genesis.CrossChainTxs, &elem)
	}

	amount, found := k.GetZetaAccounting(ctx)
	if found {
		genesis.ZetaAccounting = amount
	}
	genesis.FinalizedInbounds = k.GetAllFinalizedInbound(ctx)

	return &genesis
}
