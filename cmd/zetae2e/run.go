package main

import (
	"context"
	"errors"
	"strings"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/app"
	zetae2econfig "github.com/zeta-chain/zetacore/cmd/zetae2e/config"
	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

const flagVerbose = "verbose"
const flagConfig = "config"

const FungibleAdminMnemonic = "snow grace federal cupboard arrive fancy gym lady uniform rotate exercise either leave alien grass" // #nosec G101 - used for testing

// NewRunCmd returns the run command
// which runs the E2E from a config file describing the tests, networks, and accounts
func NewRunCmd() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "run [testname1]:[arg1],[arg2] [testname2]:[arg1],[arg2]...",
		Short: "Run one or more E2E tests with optional arguments",
		Long: `Run one or more E2E tests specified by their names and optional arguments.
For example: zetae2e run deposit:1000 withdraw: --config config.yml`,
		RunE: runE2ETest,
		Args: cobra.MinimumNArgs(1), // Ensures at least one test is provided
	}

	cmd.Flags().StringVarP(&configFile, flagConfig, "c", "", "path to the configuration file")
	err := cmd.MarkFlagRequired(flagConfig)
	if err != nil {
		panic(err)
	}

	// Retain the verbose flag
	cmd.Flags().Bool(
		flagVerbose,
		false,
		"set to true to enable verbose logging",
	)

	return cmd
}

func runE2ETest(cmd *cobra.Command, args []string) error {
	// read the config file
	configPath, err := cmd.Flags().GetString(flagConfig)
	if err != nil {
		return err
	}
	conf, err := config.ReadConfig(configPath)
	if err != nil {
		return err
	}

	// read flag
	verbose, err := cmd.Flags().GetBool(flagVerbose)
	if err != nil {
		return err
	}

	// initialize logger
	logger := runner.NewLogger(verbose, color.FgHiCyan, "e2e")

	// set config
	app.SetConfig()

	// initialize context
	ctx, cancel := context.WithCancel(context.Background())

	// get EVM address from config
	evmAddr := conf.Accounts.EVMAddress
	if !ethcommon.IsHexAddress(evmAddr) {
		cancel()
		return errors.New("invalid EVM address")
	}

	// initialize deployer runner with config
	testRunner, err := zetae2econfig.RunnerFromConfig(
		ctx,
		"e2e",
		cancel,
		conf,
		ethcommon.HexToAddress(evmAddr),
		conf.Accounts.EVMPrivKey,
		utils.FungibleAdminName, // placeholder value, not used
		FungibleAdminMnemonic,   // placeholder value, not used
		logger,
	)
	if err != nil {
		cancel()
		return err
	}

	testStartTime := time.Now()
	logger.Print("starting tests")

	// fetch the TSS address
	if err := testRunner.SetTSSAddresses(); err != nil {
		return err
	}

	// set timeout
	testRunner.CctxTimeout = 60 * time.Minute
	testRunner.ReceiptTimeout = 60 * time.Minute

	balancesBefore, err := testRunner.GetAccountBalances(true)
	if err != nil {
		cancel()
		return err
	}

	// parse test names and arguments from cmd args and run them
	userTestsConfigs, err := parseCmdArgsToE2ETestRunConfig(args)
	if err != nil {
		cancel()
		return err
	}

	testsToRun, err := testRunner.GetE2ETestsToRunByConfig(e2etests.AllE2ETests, userTestsConfigs)
	if err != nil {
		cancel()
		return err
	}
	reports, err := testRunner.RunE2ETestsIntoReport(testsToRun)
	if err != nil {
		cancel()
		return err
	}

	balancesAfter, err := testRunner.GetAccountBalances(true)
	if err != nil {
		cancel()
		return err
	}

	// Print tests completion info
	logger.Print("tests finished successfully in %s", time.Since(testStartTime).String())
	testRunner.Logger.SetColor(color.FgHiRed)
	testRunner.PrintTotalDiff(runner.GetAccountBalancesDiff(balancesBefore, balancesAfter))
	testRunner.Logger.SetColor(color.FgHiGreen)
	testRunner.PrintTestReports(reports)

	return nil
}

// parseCmdArgsToE2ETests parses command-line arguments into a slice of E2ETestRunConfig structs.
func parseCmdArgsToE2ETestRunConfig(args []string) ([]runner.E2ETestRunConfig, error) {
	tests := []runner.E2ETestRunConfig{}
	for _, arg := range args {
		parts := strings.SplitN(arg, ":", 2)
		if len(parts) != 2 {
			return nil, errors.New("command arguments should be in format: testName:testArgs")
		}
		if parts[0] == "" {
			return nil, errors.New("missing testName")
		}
		testName := parts[0]
		testArgs := []string{}
		if parts[1] != "" {
			testArgs = strings.Split(parts[1], ",")
		}
		tests = append(tests, runner.E2ETestRunConfig{
			Name: testName,
			Args: testArgs,
		})
	}
	return tests, nil
}
