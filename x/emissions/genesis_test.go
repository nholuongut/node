package emissions_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/emissions"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func TestGenesis(t *testing.T) {
	params := types.DefaultParams()
	params.ObserverSlashAmount = sdk.Int{}

	genesisState := types.GenesisState{
		Params: params,
		WithdrawableEmissions: []types.WithdrawableEmissions{
			sample.WithdrawableEmissions(t),
			sample.WithdrawableEmissions(t),
			sample.WithdrawableEmissions(t),
		},
	}

	// Init and export
	k, ctx, _, _ := keepertest.EmissionsKeeper(t)
	emissions.InitGenesis(ctx, *k, genesisState)
	got := emissions.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	// Compare genesis after init and export
	nullify.Fill(&genesisState)
	nullify.Fill(got)
	require.Equal(t, genesisState, *got)
}
