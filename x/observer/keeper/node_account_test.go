package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// Keeper Tests

func createNNodeAccount(keeper *Keeper, ctx sdk.Context, n int) []types.NodeAccount {
	items := make([]types.NodeAccount, n)
	for i := range items {
		items[i].Operator = fmt.Sprintf("%d", i)
		keeper.SetNodeAccount(ctx, items[i])
	}
	return items
}

func TestNodeAccountGet(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	items := createNNodeAccount(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetNodeAccount(ctx, item.Operator)
		require.True(t, found)
		require.Equal(t, item, rst)
	}
}
func TestNodeAccountRemove(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	items := createNNodeAccount(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveNodeAccount(ctx, item.Operator)
		_, found := keeper.GetNodeAccount(ctx, item.Operator)
		require.False(t, found)
	}
}

func TestNodeAccountGetAll(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	items := createNNodeAccount(keeper, ctx, 10)
	require.Equal(t, items, keeper.GetAllNodeAccount(ctx))
}
