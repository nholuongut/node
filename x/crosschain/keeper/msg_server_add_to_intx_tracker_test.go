package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func setupVerificationParams(zk keepertest.ZetaKeepers, ctx sdk.Context, tx_index int64, chainID int64, header ethtypes.Header, headerRLP []byte, block *ethtypes.Block) {
	params := zk.ObserverKeeper.GetParams(ctx)
	zk.ObserverKeeper.SetParams(ctx, params)
	zk.ObserverKeeper.SetBlockHeader(ctx, common.BlockHeader{
		Height:     block.Number().Int64(),
		Hash:       block.Hash().Bytes(),
		ParentHash: header.ParentHash.Bytes(),
		ChainId:    chainID,
		Header:     common.NewEthereumHeader(headerRLP),
	})
	zk.ObserverKeeper.SetChainParamsList(ctx, observertypes.ChainParamsList{ChainParams: []*observertypes.ChainParams{
		{
			ChainId:                  chainID,
			ConnectorContractAddress: block.Transactions()[tx_index].To().Hex(),
			BallotThreshold:          sdk.OneDec(),
			MinObserverDelegation:    sdk.OneDec(),
			IsSupported:              true,
		},
	}})
	zk.ObserverKeeper.SetCrosschainFlags(ctx, observertypes.CrosschainFlags{
		BlockHeaderVerificationFlags: &observertypes.BlockHeaderVerificationFlags{
			IsEthTypeChainEnabled: true,
			IsBtcTypeChainEnabled: false,
		},
	})
}

func TestMsgServer_AddToInTxTracker(t *testing.T) {
	t.Run("fail normal user submit without proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		tx_hash := "string"

		chainID := getValidEthChainID(t)
		setSupportedChain(ctx, zk, chainID)

		msgServer := keeper.NewMsgServerImpl(*k)

		_, err := msgServer.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID,
			TxHash:    tx_hash,
			CoinType:  common.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		})
		require.ErrorIs(t, err, observertypes.ErrNotAuthorized)
		_, found := k.GetInTxTracker(ctx, chainID, tx_hash)
		require.False(t, found)
	})

	t.Run("admin add  tx tracker", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)

		tx_hash := "string"
		chainID := getValidEthChainID(t)
		setSupportedChain(ctx, zk, chainID)

		msgServer := keeper.NewMsgServerImpl(*k)

		_, err := msgServer.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    tx_hash,
			CoinType:  common.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		})
		require.NoError(t, err)
		_, found := k.GetInTxTracker(ctx, chainID, tx_hash)
		require.True(t, found)
	})

	t.Run("admin submit fake tracker", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)

		tx_hash := "string"
		chainID := getValidEthChainID(t)
		setSupportedChain(ctx, zk, chainID)

		msgServer := keeper.NewMsgServerImpl(*k)

		_, err := msgServer.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    "Malicious TX HASH",
			CoinType:  common.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		})
		require.NoError(t, err)
		_, found := k.GetInTxTracker(ctx, chainID, "Malicious TX HASH")
		require.True(t, found)
		_, found = k.GetInTxTracker(ctx, chainID, tx_hash)
		require.False(t, found)
	})

	// Commented out as these tests don't work without using RPC
	// TODO: Reenable these tests
	// https://github.com/zeta-chain/node/issues/1875
	//t.Run("add proof based tracker with correct proof", func(t *testing.T) {
	//	k, ctx, _, zk := keepertest.CrosschainKeeper(t)
	//
	//	chainID := int64(5)
	//
	//	txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
	//	require.NoError(t, err)
	//	setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
	//	msgServer := keeper.NewMsgServerImpl(*k)
	//
	//	_, err = msgServer.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
	//		Creator:   sample.AccAddress(),
	//		ChainId:   chainID,
	//		TxHash:    tx.Hash().Hex(),
	//		CoinType:  common.CoinType_Zeta,
	//		Proof:     proof,
	//		BlockHash: block.Hash().Hex(),
	//		TxIndex:   txIndex,
	//	})
	//	require.NoError(t, err)
	//	_, found := k.GetInTxTracker(ctx, chainID, tx.Hash().Hex())
	//	require.True(t, found)
	//})
	//t.Run("fail to add proof based tracker with wrong tx hash", func(t *testing.T) {
	//	k, ctx, _, zk := keepertest.CrosschainKeeper(t)
	//
	//	chainID := getValidEthChainID(t)
	//
	//	txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
	//	require.NoError(t, err)
	//	setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
	//	msgServer := keeper.NewMsgServerImpl(*k)
	//
	//	_, err = msgServer.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
	//		Creator:   sample.AccAddress(),
	//		ChainId:   chainID,
	//		TxHash:    "fake_hash",
	//		CoinType:  common.CoinType_Zeta,
	//		Proof:     proof,
	//		BlockHash: block.Hash().Hex(),
	//		TxIndex:   txIndex,
	//	})
	//	require.ErrorIs(t, err, types.ErrTxBodyVerificationFail)
	//	_, found := k.GetInTxTracker(ctx, chainID, tx.Hash().Hex())
	//	require.False(t, found)
	//})
	//t.Run("fail to add proof based tracker with wrong chain id", func(t *testing.T) {
	//	k, ctx, _, zk := keepertest.CrosschainKeeper(t)
	//
	//	chainID := getValidEthChainID(t)
	//
	//	txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
	//	require.NoError(t, err)
	//	setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
	//
	//	msgServer := keeper.NewMsgServerImpl(*k)
	//
	//	_, err = msgServer.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
	//		Creator:   sample.AccAddress(),
	//		ChainId:   97,
	//		TxHash:    tx.Hash().Hex(),
	//		CoinType:  common.CoinType_Zeta,
	//		Proof:     proof,
	//		BlockHash: block.Hash().Hex(),
	//		TxIndex:   txIndex,
	//	})
	//	require.ErrorIs(t, err, observertypes.ErrSupportedChains)
	//	_, found := k.GetInTxTracker(ctx, chainID, tx.Hash().Hex())
	//	require.False(t, found)
	//})
}
