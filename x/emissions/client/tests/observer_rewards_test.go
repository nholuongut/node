package querytests

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	emissionscli "github.com/zeta-chain/zetacore/x/emissions/client/cli"
	emissionskeeper "github.com/zeta-chain/zetacore/x/emissions/keeper"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	observercli "github.com/zeta-chain/zetacore/x/observer/client/cli"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (s *CliTestSuite) TestObserverRewards() {
	emissionPool := "800000000000000000000azeta"
	val := s.network.Validators[0]

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, emissionscli.CmdListPoolAddresses(), []string{"--output", "json"})
	s.Require().NoError(err)
	resPools := emissionstypes.QueryListPoolAddressesResponse{}
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resPools))
	txArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewInt(10))).String()),
	}

	// Fund the emission pool to start the emission process
	sendArgs := []string{val.Address.String(),
		resPools.EmissionModuleAddress, emissionPool}
	args := append(sendArgs, txArgs...)
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewSendTxCmd(), args)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	// Collect parameter values and build assertion map for the randomised ballot set created
	emissionFactors := emissionstypes.QueryGetEmissionsFactorsResponse{}
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, emissionscli.CmdGetEmmisonsFactors(), []string{"--output", "json"})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &emissionFactors))
	emissionParams := emissionstypes.QueryParamsResponse{}
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, emissionscli.CmdQueryParams(), []string{"--output", "json"})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &emissionParams))
	observerParams := observertypes.QueryParamsResponse{}
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, observercli.CmdQueryParams(), []string{"--output", "json"})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &observerParams))
	_, err = s.network.WaitForHeight(s.ballots[0].BallotCreationHeight + observerParams.Params.BallotMaturityBlocks)
	s.Require().NoError(err)
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, emissionscli.CmdGetEmmisonsFactors(), []string{"--output", "json"})
	resFactorsNewBlocks := emissionstypes.QueryGetEmissionsFactorsResponse{}
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resFactorsNewBlocks))
	// Duration factor is calculated in the same block,so we need to query based from the committed state at which the distribution is done
	// Would be cleaner to use `--height` flag, but it is not supported by the ExecTestCLICmd function yet
	emissionFactors.DurationFactor = resFactorsNewBlocks.DurationFactor
	asertValues := CalculateObserverRewards(&s.Suite, s.ballots, emissionParams.Params.ObserverEmissionPercentage, emissionFactors.ReservesFactor, emissionFactors.BondFactor, emissionFactors.DurationFactor)

	// Assert withdrawable rewards for each validator
	resAvailable := emissionstypes.QueryShowAvailableEmissionsResponse{}
	for i := 0; i < len(s.network.Validators); i++ {
		out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, emissionscli.CmdShowAvailableEmissions(), []string{s.network.Validators[i].Address.String(), "--output", "json"})
		s.Require().NoError(err)
		s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resAvailable))
		s.Require().Equal(sdk.NewCoin(config.BaseDenom, asertValues[s.network.Validators[i].Address.String()]).String(), resAvailable.Amount, "Validator %s has incorrect withdrawable rewards", s.network.Validators[i].Address.String())
	}

}

func CalculateObserverRewards(s *suite.Suite, ballots []*observertypes.Ballot, observerEmissionPercentage, reservesFactor, bondFactor, durationFactor string) map[string]sdkmath.Int {
	calculatedDistributer := map[string]sdkmath.Int{}
	//blockRewards := sdk.MustNewDecFromStr(reservesFactor).Mul(sdk.MustNewDecFromStr(bondFactor)).Mul(sdk.MustNewDecFromStr(durationFactor))
	blockRewards, err := emissionskeeper.CalculateFixedValidatorRewards(emissionstypes.AvgBlockTime)
	s.Require().NoError(err)
	observerRewards := sdk.MustNewDecFromStr(observerEmissionPercentage).Mul(blockRewards).TruncateInt()
	rewardsDistributer := map[string]int64{}
	totalRewardsUnits := int64(0)
	// BuildRewardsDistribution has a separate unit test
	for _, ballot := range ballots {
		totalRewardsUnits = totalRewardsUnits + ballot.BuildRewardsDistribution(rewardsDistributer)
	}
	rewardPerUnit := observerRewards.Quo(sdk.NewInt(totalRewardsUnits))
	for address, units := range rewardsDistributer {
		if units == 0 {
			calculatedDistributer[address] = sdk.ZeroInt()
			continue
		}
		if units < 0 {
			calculatedDistributer[address] = sdk.ZeroInt()
			continue
		}
		if units > 0 {
			calculatedDistributer[address] = rewardPerUnit.Mul(sdkmath.NewInt(units))
		}
	}
	return calculatedDistributer
}
