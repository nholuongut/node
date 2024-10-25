package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_UpdateZRC20PausedStatus(t *testing.T) {
	t.Run("can update the paused status of zrc20", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		assertUnpaused := func(zrc20 string) {
			fc, found := k.GetForeignCoins(ctx, zrc20)
			require.True(t, found)
			require.False(t, fc.Paused)
		}
		assertPaused := func(zrc20 string) {
			fc, found := k.GetForeignCoins(ctx, zrc20)
			require.True(t, found)
			require.True(t, fc.Paused)
		}

		// setup zrc20
		zrc20A, zrc20B, zrc20C := sample.EthAddress().String(), sample.EthAddress().String(), sample.EthAddress().String()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20A))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20B))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20C))
		assertUnpaused(zrc20A)
		assertUnpaused(zrc20B)
		assertUnpaused(zrc20C)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)

		// can pause zrc20
		_, err := msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20A,
				zrc20B,
			},
			types.UpdatePausedStatusAction_PAUSE,
		))
		require.NoError(t, err)
		assertPaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		// can unpause zrc20
		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20A,
			},
			types.UpdatePausedStatusAction_UNPAUSE,
		))
		require.NoError(t, err)
		assertUnpaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)

		// can pause already paused zrc20
		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20B,
			},
			types.UpdatePausedStatusAction_PAUSE,
		))
		require.NoError(t, err)
		assertUnpaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		// can unpause already unpaused zrc20
		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20C,
			},
			types.UpdatePausedStatusAction_UNPAUSE,
		))
		require.NoError(t, err)
		assertUnpaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)

		// can pause all zrc20
		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20A,
				zrc20B,
				zrc20C,
			},
			types.UpdatePausedStatusAction_PAUSE,
		))
		require.NoError(t, err)
		assertPaused(zrc20A)
		assertPaused(zrc20B)
		assertPaused(zrc20C)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		// can unpause all zrc20
		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20A,
				zrc20B,
				zrc20C,
			},
			types.UpdatePausedStatusAction_UNPAUSE,
		))
		require.NoError(t, err)
		assertUnpaused(zrc20A)
		assertUnpaused(zrc20B)
		assertUnpaused(zrc20C)
	})

	t.Run("should fail if invalid message", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		invalidMsg := types.NewMsgUpdateZRC20PausedStatus(admin, []string{}, types.UpdatePausedStatusAction_PAUSE)
		require.ErrorIs(t, invalidMsg.ValidateBasic(), sdkerrors.ErrInvalidRequest)

		_, err := msgServer.UpdateZRC20PausedStatus(ctx, invalidMsg)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})

	t.Run("should fail if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, false)

		_, err := msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{sample.EthAddress().String()},
			types.UpdatePausedStatusAction_PAUSE,
		))
		require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, false)

		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{sample.EthAddress().String()},
			types.UpdatePausedStatusAction_UNPAUSE,
		))

		require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("should fail if zrc20 does not exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)

		zrc20A, zrc20B := sample.EthAddress().String(), sample.EthAddress().String()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20A))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20B))

		_, err := msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20A,
				sample.EthAddress().String(),
				zrc20B,
			},
			types.UpdatePausedStatusAction_PAUSE,
		))
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})
}
