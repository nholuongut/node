package keeper_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_RefundAmountOnZetaChainGas(t *testing.T) {
	t.Run("should refund amount zrc20 gas on zeta chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		err := k.RefundAmountOnZetaChainGas(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
			OutboundTxParams: []*types.OutboundTxParams{{
				Amount: math.NewUint(42),
			}},
		},
			sender,
		)
		require.NoError(t, err)
		balance, err := zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20, sender)
		require.NoError(t, err)
		require.Equal(t, uint64(42), balance.Uint64())
	})
	t.Run("should refund inbound amount zrc20 gas on zeta chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		err := k.RefundAmountOnZetaChainGas(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
		},
			sender,
		)
		require.NoError(t, err)
		balance, err := zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20, sender)
		require.NoError(t, err)
		require.Equal(t, uint64(20), balance.Uint64())
	})
	t.Run("failed refund zrc20 gas on zeta chain if gas coin not found", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		err := k.RefundAmountOnZetaChainGas(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
			OutboundTxParams: []*types.OutboundTxParams{{
				Amount: math.NewUint(42),
			}},
		},

			sender,
		)
		require.ErrorContains(t, err, types.ErrForeignCoinNotFound.Error())
	})
	t.Run("failed refund amount zrc20 gas on zeta chain if amount is 0", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		err := k.RefundAmountOnZetaChainGas(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.ZeroUint(),
			},
			OutboundTxParams: []*types.OutboundTxParams{{
				Amount: math.ZeroUint(),
			}},
		},
			sender,
		)
		require.ErrorContains(t, err, "no amount to refund")
	})

}

func TestKeeper_RefundAmountOnZetaChainZeta(t *testing.T) {
	t.Run("should refund amount on zeta chain", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID(t)

		err := k.RefundAmountOnZetaChainZeta(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
			OutboundTxParams: []*types.OutboundTxParams{{
				Amount: math.NewUint(42),
			}},
		},
			sender,
		)
		require.NoError(t, err)
		coin := sdkk.BankKeeper.GetBalance(ctx, sdk.AccAddress(sender.Bytes()), config.BaseDenom)
		fmt.Println(coin.Amount.String())
		require.Equal(t, "42", coin.Amount.String())
	})
	t.Run("should refund inbound amount on zeta chain if outbound is not present", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID(t)

		err := k.RefundAmountOnZetaChainZeta(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
		},
			sender,
		)
		require.NoError(t, err)
		coin := sdkk.BankKeeper.GetBalance(ctx, sdk.AccAddress(sender.Bytes()), config.BaseDenom)
		require.Equal(t, "20", coin.Amount.String())
	})
	t.Run("failed refund amount on zeta chain amount is 0", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID(t)

		err := k.RefundAmountOnZetaChainZeta(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.ZeroUint(),
			},
			OutboundTxParams: []*types.OutboundTxParams{{
				Amount: math.ZeroUint(),
			}},
		},
			sender,
		)
		require.ErrorContains(t, err, "no amount to refund")
	})
}

func TestKeeper_RefundAmountOnZetaChainERC20(t *testing.T) {
	t.Run("should refund amount on zeta chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		asset := sample.EthAddress().String()
		sender := sample.EthAddress()
		chainID := getValidEthChainID(t)

		// deploy zrc20
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
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

		err := k.RefundAmountOnZetaChainERC20(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: chainID,
				Sender:        sender.String(),
				Asset:         asset,
				Amount:        math.NewUint(42),
			},
			OutboundTxParams: []*types.OutboundTxParams{{
				Amount: math.NewUint(42),
			}},
		},
			sender,
		)
		require.NoError(t, err)

		// check amount deposited in balance
		balance, err := zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20Addr, sender)
		require.NoError(t, err)
		require.Equal(t, uint64(42), balance.Uint64())

		// can refund again
		err = k.RefundAmountOnZetaChainERC20(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: chainID,
				Sender:        sender.String(),
				Asset:         asset,
				Amount:        math.NewUint(42),
			}},
			sender,
		)
		require.NoError(t, err)
		balance, err = zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20Addr, sender)
		require.NoError(t, err)
		require.Equal(t, uint64(84), balance.Uint64())
	})

	t.Run("should fail with invalid cctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		err := k.RefundAmountOnZetaChainERC20(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: common.CoinType_Zeta,
				Amount:   math.NewUint(42),
			}},
			sample.EthAddress(),
		)
		require.ErrorContains(t, err, "unsupported coin type")

		err = k.RefundAmountOnZetaChainERC20(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: common.CoinType_Gas,
			}},
			sample.EthAddress(),
		)
		require.ErrorContains(t, err, "unsupported coin type")

		err = k.RefundAmountOnZetaChainERC20(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: 999999,
				Amount:        math.NewUint(42),
			}},
			sample.EthAddress(),
		)
		require.ErrorContains(t, err, "only EVM chains are supported")

		err = k.RefundAmountOnZetaChainERC20(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: getValidEthChainID(t),
				Sender:        sample.EthAddress().String(),
				Amount:        math.Uint{},
			}},
			sample.EthAddress(),
		)
		require.ErrorContains(t, err, "no amount to refund")

		err = k.RefundAmountOnZetaChainERC20(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: getValidEthChainID(t),
				Sender:        sample.EthAddress().String(),
				Amount:        math.ZeroUint(),
			}},
			sample.EthAddress(),
		)
		require.ErrorContains(t, err, "no amount to refund")

		// the foreign coin has not been set
		err = k.RefundAmountOnZetaChainERC20(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: getValidEthChainID(t),
				Sender:        sample.EthAddress().String(),
				Asset:         sample.EthAddress().String(),
				Amount:        math.NewUint(42),
			}},
			sample.EthAddress(),
		)
		require.ErrorContains(t, err, "zrc not found")
	})
}
