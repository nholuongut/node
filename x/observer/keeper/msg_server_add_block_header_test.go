package keeper_test

import (
	"encoding/json"
	"math/rand"
	"os"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_AddBlockHeader(t *testing.T) {
	header, header2, header3, err := ethHeaders()
	require.NoError(t, err)
	header1RLP, err := rlp.EncodeToBytes(header)
	require.NoError(t, err)
	header2RLP, err := rlp.EncodeToBytes(header2)
	require.NoError(t, err)
	header3RLP, err := rlp.EncodeToBytes(header3)
	require.NoError(t, err)

	r := rand.New(rand.NewSource(9))
	validator := sample.Validator(t, r)
	observerAddress, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
	require.NoError(t, err)
	// Add tests for btc headers : https://github.com/zeta-chain/node/issues/1336
	tt := []struct {
		name                  string
		msg                   *types.MsgAddBlockHeader
		IsEthTypeChainEnabled bool
		IsBtcTypeChainEnabled bool
		validator             stakingtypes.Validator
		wantErr               require.ErrorAssertionFunc
	}{
		{
			name: "success submit eth header",
			msg: &types.MsgAddBlockHeader{
				Creator:   observerAddress.String(),
				ChainId:   common.GoerliLocalnetChain().ChainId,
				BlockHash: header.Hash().Bytes(),
				Height:    1,
				Header:    common.NewEthereumHeader(header1RLP),
			},
			IsEthTypeChainEnabled: true,
			IsBtcTypeChainEnabled: true,
			validator:             validator,
			wantErr:               require.NoError,
		},
		{
			name: "failure submit eth header eth disabled",
			msg: &types.MsgAddBlockHeader{
				Creator:   observerAddress.String(),
				ChainId:   common.GoerliLocalnetChain().ChainId,
				BlockHash: header.Hash().Bytes(),
				Height:    1,
				Header:    common.NewEthereumHeader(header1RLP),
			},
			IsEthTypeChainEnabled: false,
			IsBtcTypeChainEnabled: true,
			validator:             validator,
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, types.ErrBlockHeaderVerificationDisabled)
			},
		},
		{
			name: "failure submit eth header eth disabled",
			msg: &types.MsgAddBlockHeader{
				Creator:   sample.AccAddress(),
				ChainId:   common.GoerliLocalnetChain().ChainId,
				BlockHash: header.Hash().Bytes(),
				Height:    1,
				Header:    common.NewEthereumHeader(header1RLP),
			},
			IsEthTypeChainEnabled: false,
			IsBtcTypeChainEnabled: true,
			validator:             validator,
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, types.ErrNotAuthorizedPolicy)
			},
		},
		{
			name: "should succeed if block header parent does exist",
			msg: &types.MsgAddBlockHeader{
				Creator:   observerAddress.String(),
				ChainId:   common.GoerliLocalnetChain().ChainId,
				BlockHash: header2.Hash().Bytes(),
				Height:    2,
				Header:    common.NewEthereumHeader(header2RLP),
			},
			IsEthTypeChainEnabled: true,
			IsBtcTypeChainEnabled: true,
			validator:             validator,
			wantErr:               require.NoError,
		},
		// These tests don't work when using the static headers, the previous sample were also not correct (header3 used to be nil)
		// The second test mention it should success but assert an error
		// TODO: fix these tests
		// https://github.com/zeta-chain/node/issues/1875
		//{
		//	name: "should fail if block header parent does not exist",
		//	msg: &types.MsgAddBlockHeader{
		//		Creator:   observerAddress.String(),
		//		ChainId:   common.GoerliLocalnetChain().ChainId,
		//		BlockHash: header3.Hash().Bytes(),
		//		Height:    3,
		//		Header:    common.NewEthereumHeader(header3RLP),
		//	},
		//	IsEthTypeChainEnabled: true,
		//	IsBtcTypeChainEnabled: true,
		//	validator:             validator,
		//	wantErr: func(t require.TestingT, err error, i ...interface{}) {
		//		require.Error(t, err)
		//	},
		//},
		//{
		//	name: "should succeed to post 3rd header if 2nd header is posted",
		//	msg: &types.MsgAddBlockHeader{
		//		Creator:   observerAddress.String(),
		//		ChainId:   common.GoerliLocalnetChain().ChainId,
		//		BlockHash: header3.Hash().Bytes(),
		//		Height:    3,
		//		Header:    common.NewEthereumHeader(header3RLP),
		//	},
		//	IsEthTypeChainEnabled: true,
		//	IsBtcTypeChainEnabled: true,
		//	validator:             validator,
		//	wantErr: func(t require.TestingT, err error, i ...interface{}) {
		//		require.Error(t, err)
		//	},
		//},
		{
			name: "should fail if chain is not supported",
			msg: &types.MsgAddBlockHeader{
				Creator:   observerAddress.String(),
				ChainId:   9999,
				BlockHash: header3.Hash().Bytes(),
				Height:    3,
				Header:    common.NewEthereumHeader(header3RLP),
			},
			IsEthTypeChainEnabled: true,
			IsBtcTypeChainEnabled: true,
			validator:             validator,
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, types.ErrSupportedChains)
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			k, ctx, _, _ := keepertest.ObserverKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			k.SetObserverSet(ctx, types.ObserverSet{
				ObserverList: []string{observerAddress.String()},
			})
			k.GetStakingKeeper().SetValidator(ctx, tc.validator)
			k.SetCrosschainFlags(ctx, types.CrosschainFlags{
				IsInboundEnabled:      true,
				IsOutboundEnabled:     true,
				GasPriceIncreaseFlags: nil,
				BlockHeaderVerificationFlags: &types.BlockHeaderVerificationFlags{
					IsEthTypeChainEnabled: tc.IsEthTypeChainEnabled,
					IsBtcTypeChainEnabled: tc.IsBtcTypeChainEnabled,
				},
			})

			setSupportedChain(ctx, *k, common.GoerliLocalnetChain().ChainId)

			_, err := srv.AddBlockHeader(ctx, tc.msg)
			tc.wantErr(t, err)
			if err == nil {
				bhs, found := k.GetBlockHeaderState(ctx, tc.msg.ChainId)
				require.True(t, found)
				require.Equal(t, tc.msg.Height, bhs.LatestHeight)
			}
		})
	}
}

func ethHeaders() (*ethtypes.Header, *ethtypes.Header, *ethtypes.Header, error) {
	header1, err := readHeader("./testdata/header_sepolia_5000000.json")
	if err != nil {
		return nil, nil, nil, err
	}
	header2, err := readHeader("./testdata/header_sepolia_5000001.json")
	if err != nil {
		return nil, nil, nil, err
	}
	header3, err := readHeader("./testdata/header_sepolia_5000002.json")
	if err != nil {
		return nil, nil, nil, err
	}
	return header1, header2, header3, nil
}

// readReceipt reads a receipt from a file.
// TODO: centralize test data
// https://github.com/zeta-chain/node/issues/1874
func readHeader(filename string) (*ethtypes.Header, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var NewHeader ethtypes.Header
	err = decoder.Decode(&NewHeader)
	return &NewHeader, err
}
