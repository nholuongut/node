package keeper_test

import (
	"encoding/hex"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func setObservers(t *testing.T, k *keeper.Keeper, ctx sdk.Context, zk keepertest.ZetaKeepers) []string {
	validators := k.GetStakingKeeper().GetAllValidators(ctx)

	validatorAddressListFormatted := make([]string, len(validators))
	for i, validator := range validators {
		valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		require.NoError(t, err)
		addressTmp, err := sdk.AccAddressFromHexUnsafe(hex.EncodeToString(valAddr.Bytes()))
		require.NoError(t, err)
		validatorAddressListFormatted[i] = addressTmp.String()
	}

	// Add validator to the observer list for voting
	zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{
		ObserverList: validatorAddressListFormatted,
	})
	return validatorAddressListFormatted
}

// TODO: Complete the test cases
// https://github.com/zeta-chain/node/issues/1542
func TestKeeper_VoteOnObservedInboundTx(t *testing.T) {
	t.Run("successfully vote on evm deposit", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		validatorList := setObservers(t, k, ctx, zk)
		to, from := int64(1337), int64(101)
		chains := zk.ObserverKeeper.GetSupportedChains(ctx)
		for _, chain := range chains {
			if common.IsEVMChain(chain.ChainId) {
				from = chain.ChainId
			}
			if common.IsZetaChain(chain.ChainId) {
				to = chain.ChainId
			}
		}
		msg := sample.InboundVote(0, from, to)
		for _, validatorAddr := range validatorList {
			msg.Creator = validatorAddr
			_, err := msgServer.VoteOnObservedInboundTx(
				ctx,
				&msg,
			)
			require.NoError(t, err)
		}
		ballot, _, _ := zk.ObserverKeeper.FindBallot(
			ctx,
			msg.Digest(),
			zk.ObserverKeeper.GetSupportedChainFromChainID(ctx, msg.SenderChainId),
			observertypes.ObservationType_InBoundTx,
		)
		require.Equal(t, ballot.BallotStatus, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		cctx, found := k.GetCrossChainTx(ctx, msg.Digest())
		require.True(t, found)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_OutboundMined)
		require.Equal(t, cctx.InboundTxParams.TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})

	t.Run("prevent double event submission", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)

		// MsgServer for the crosschain keeper
		msgServer := keeper.NewMsgServerImpl(*k)

		// Set the chain ids we want to use to be valid
		params := observertypes.DefaultParams()
		zk.ObserverKeeper.SetParams(
			ctx, params,
		)

		// Convert the validator address into a user address.
		validators := k.GetStakingKeeper().GetAllValidators(ctx)
		validatorAddress := validators[0].OperatorAddress
		valAddr, _ := sdk.ValAddressFromBech32(validatorAddress)
		addresstmp, _ := sdk.AccAddressFromHexUnsafe(hex.EncodeToString(valAddr.Bytes()))
		validatorAddr := addresstmp.String()

		// Add validator to the observer list for voting
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{
			ObserverList: []string{validatorAddr},
		})

		// Vote on the FIRST message.
		msg := &types.MsgVoteOnObservedInboundTx{
			Creator:       validatorAddr,
			Sender:        "0x954598965C2aCdA2885B037561526260764095B8",
			SenderChainId: 1337, // ETH
			Receiver:      "0x954598965C2aCdA2885B037561526260764095B8",
			ReceiverChain: 101, // zetachain
			Amount:        sdkmath.NewUintFromString("10000000"),
			Message:       "",
			InBlockHeight: 1,
			GasLimit:      1000000000,
			InTxHash:      "0x7a900ef978743f91f57ca47c6d1a1add75df4d3531da17671e9cf149e1aefe0b",
			CoinType:      0, // zeta
			TxOrigin:      "0x954598965C2aCdA2885B037561526260764095B8",
			Asset:         "",
			EventIndex:    1,
		}
		_, err := msgServer.VoteOnObservedInboundTx(
			ctx,
			msg,
		)
		require.NoError(t, err)

		// Check that the vote passed
		ballot, found := zk.ObserverKeeper.GetBallot(ctx, msg.Digest())
		require.True(t, found)
		require.Equal(t, ballot.BallotStatus, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		//Perform the SAME event. Except, this time, we resubmit the event.
		msg2 := &types.MsgVoteOnObservedInboundTx{
			Creator:       validatorAddr,
			Sender:        "0x954598965C2aCdA2885B037561526260764095B8",
			SenderChainId: 1337,
			Receiver:      "0x954598965C2aCdA2885B037561526260764095B8",
			ReceiverChain: 101,
			Amount:        sdkmath.NewUintFromString("10000000"),
			Message:       "",
			InBlockHeight: 1,
			GasLimit:      1000000001, // <---- Change here
			InTxHash:      "0x7a900ef978743f91f57ca47c6d1a1add75df4d3531da17671e9cf149e1aefe0b",
			CoinType:      0,
			TxOrigin:      "0x954598965C2aCdA2885B037561526260764095B8",
			Asset:         "",
			EventIndex:    1,
		}

		_, err = msgServer.VoteOnObservedInboundTx(
			ctx,
			msg2,
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrObservedTxAlreadyFinalized)
		_, found = zk.ObserverKeeper.GetBallot(ctx, msg2.Digest())
		require.False(t, found)
	})
}

func TestStatus_ChangeStatus(t *testing.T) {
	tt := []struct {
		Name         string
		Status       types.Status
		NonErrStatus types.CctxStatus
		Msg          string
		IsErr        bool
		ErrStatus    types.CctxStatus
	}{
		{
			Name: "Transition on finalize Inbound",
			Status: types.Status{
				Status:              types.CctxStatus_PendingInbound,
				StatusMessage:       "Getting InTX Votes",
				LastUpdateTimestamp: 0,
			},
			Msg:          "Got super majority and finalized Inbound",
			NonErrStatus: types.CctxStatus_PendingOutbound,
			ErrStatus:    types.CctxStatus_Aborted,
			IsErr:        false,
		},
		{
			Name: "Transition on finalize Inbound Fail",
			Status: types.Status{
				Status:              types.CctxStatus_PendingInbound,
				StatusMessage:       "Getting InTX Votes",
				LastUpdateTimestamp: 0,
			},
			Msg:          "Got super majority and finalized Inbound",
			NonErrStatus: types.CctxStatus_OutboundMined,
			ErrStatus:    types.CctxStatus_Aborted,
			IsErr:        false,
		},
	}
	for _, test := range tt {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			test.Status.ChangeStatus(test.NonErrStatus, test.Msg)
			if test.IsErr {
				require.Equal(t, test.ErrStatus, test.Status.Status)
			} else {
				require.Equal(t, test.NonErrStatus, test.Status.Status)
			}
		})
	}
}
