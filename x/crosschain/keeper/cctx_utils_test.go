package keeper_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestGetRevertGasLimit(t *testing.T) {
	t.Run("should return 0 if no inbound tx params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		gasLimit, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{})
		require.NoError(t, err)
		require.Equal(t, uint64(0), gasLimit)
	})

	t.Run("should return 0 if coin type is not gas or erc20", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		gasLimit, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: common.CoinType_Zeta,
			}})
		require.NoError(t, err)
		require.Equal(t, uint64(0), gasLimit)
	})

	t.Run("should return the gas limit of the gas token", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		gas := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foo", "FOO")

		_, err := zk.FungibleKeeper.UpdateZRC20GasLimit(ctx, gas, big.NewInt(42))
		require.NoError(t, err)

		gasLimit, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: chainID,
			}})
		require.NoError(t, err)
		require.Equal(t, uint64(42), gasLimit)
	})

	t.Run("should return the gas limit of the associated asset", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		asset := sample.EthAddress().String()
		zrc20Addr := deployZRC20(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			chainID,
			"bar",
			asset,
			"bar",
		)

		_, err := zk.FungibleKeeper.UpdateZRC20GasLimit(ctx, zrc20Addr, big.NewInt(42))
		require.NoError(t, err)

		gasLimit, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: chainID,
				Asset:         asset,
			}})
		require.NoError(t, err)
		require.Equal(t, uint64(42), gasLimit)
	})

	t.Run("should fail if no gas coin found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		_, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: 999999,
			}})
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})

	t.Run("should fail if query gas limit for gas coin fails", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)

		zk.FungibleKeeper.SetForeignCoins(ctx, fungibletypes.ForeignCoins{
			Zrc20ContractAddress: sample.EthAddress().String(),
			ForeignChainId:       chainID,
			CoinType:             common.CoinType_Gas,
		})

		// no contract deployed therefore will fail
		_, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: chainID,
			}})
		require.ErrorIs(t, err, fungibletypes.ErrContractCall)
	})

	t.Run("should fail if no asset found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		_, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: 999999,
			}})
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})

	t.Run("should fail if query gas limit for asset fails", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)
		asset := sample.EthAddress().String()

		zk.FungibleKeeper.SetForeignCoins(ctx, fungibletypes.ForeignCoins{
			Zrc20ContractAddress: sample.EthAddress().String(),
			ForeignChainId:       chainID,
			CoinType:             common.CoinType_ERC20,
			Asset:                asset,
		})

		// no contract deployed therefore will fail
		_, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: chainID,
				Asset:         asset,
			}})
		require.ErrorIs(t, err, fungibletypes.ErrContractCall)
	})
}

func TestGetAbortedAmount(t *testing.T) {
	amount := sdkmath.NewUint(100)
	t.Run("should return the inbound amount if outbound not present", func(t *testing.T) {
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				Amount: amount,
			},
		}
		a := crosschainkeeper.GetAbortedAmount(cctx)
		require.Equal(t, amount, a)
	})
	t.Run("should return the amount outbound amount", func(t *testing.T) {
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				Amount: sdkmath.ZeroUint(),
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{Amount: amount},
			},
		}
		a := crosschainkeeper.GetAbortedAmount(cctx)
		require.Equal(t, amount, a)
	})
	t.Run("should return the zero if outbound amount is not present and inbound is 0", func(t *testing.T) {
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				Amount: sdkmath.ZeroUint(),
			},
		}
		a := crosschainkeeper.GetAbortedAmount(cctx)
		require.Equal(t, sdkmath.ZeroUint(), a)
	})
	t.Run("should return the zero if no amounts are present", func(t *testing.T) {
		cctx := types.CrossChainTx{}
		a := crosschainkeeper.GetAbortedAmount(cctx)
		require.Equal(t, sdkmath.ZeroUint(), a)
	})
}
