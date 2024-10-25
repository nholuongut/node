package common_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
)

func Test_GetAzetaDecFromAmountInZeta(t *testing.T) {
	tt := []struct {
		name        string
		zetaAmount  string
		err         require.ErrorAssertionFunc
		azetaAmount sdk.Dec
	}{
		{
			name:        "valid zeta amount",
			zetaAmount:  "210000000",
			err:         require.NoError,
			azetaAmount: sdk.MustNewDecFromStr("210000000000000000000000000"),
		},
		{
			name:        "very high zeta amount",
			zetaAmount:  "21000000000000000000",
			err:         require.NoError,
			azetaAmount: sdk.MustNewDecFromStr("21000000000000000000000000000000000000"),
		},
		{
			name:        "very low zeta amount",
			zetaAmount:  "1",
			err:         require.NoError,
			azetaAmount: sdk.MustNewDecFromStr("1000000000000000000"),
		},
		{
			name:        "zero zeta amount",
			zetaAmount:  "0",
			err:         require.NoError,
			azetaAmount: sdk.MustNewDecFromStr("0"),
		},
		{
			name:        "decimal zeta amount",
			zetaAmount:  "0.1",
			err:         require.NoError,
			azetaAmount: sdk.MustNewDecFromStr("100000000000000000"),
		},
		{
			name:        "invalid zeta amount",
			zetaAmount:  "%%%%%$#",
			err:         require.Error,
			azetaAmount: sdk.MustNewDecFromStr("0"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			azeta, err := common.GetAzetaDecFromAmountInZeta(tc.zetaAmount)
			tc.err(t, err)
			if err == nil {
				require.Equal(t, tc.azetaAmount, azeta)
			}
		})
	}

}
