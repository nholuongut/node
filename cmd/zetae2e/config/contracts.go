package config

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	zetaeth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zeta.eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/connectorzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/wzeta.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/contracts/contextapp"
	"github.com/zeta-chain/zetacore/e2e/contracts/erc20"
	"github.com/zeta-chain/zetacore/e2e/contracts/zevmswap"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// setContractsFromConfigs get EVM contracts from config
func setContractsFromConfig(r *runner.E2ERunner, conf config.Config) error {
	var err error

	// set EVM contracts
	if c := conf.Contracts.EVM.ZetaEthAddress; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid ZetaEthAddress: %s", c)
		}
		r.ZetaEthAddr = ethcommon.HexToAddress(c)
		r.ZetaEth, err = zetaeth.NewZetaEth(r.ZetaEthAddr, r.EVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.EVM.ConnectorEthAddr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid ConnectorEthAddr: %s", c)
		}
		r.ConnectorEthAddr = ethcommon.HexToAddress(c)
		r.ConnectorEth, err = zetaconnectoreth.NewZetaConnectorEth(r.ConnectorEthAddr, r.EVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.EVM.CustodyAddr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid CustodyAddr: %s", c)
		}
		r.ERC20CustodyAddr = ethcommon.HexToAddress(c)
		r.ERC20Custody, err = erc20custody.NewERC20Custody(r.ERC20CustodyAddr, r.EVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.EVM.ERC20; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid ERC20: %s", c)
		}
		r.ERC20Addr = ethcommon.HexToAddress(c)
		r.ERC20, err = erc20.NewERC20(r.ERC20Addr, r.EVMClient)
		if err != nil {
			return err
		}
	}

	// set Zevm contracts
	if c := conf.Contracts.ZEVM.SystemContractAddr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid SystemContractAddr: %s", c)
		}
		r.SystemContractAddr = ethcommon.HexToAddress(c)
		r.SystemContract, err = systemcontract.NewSystemContract(r.SystemContractAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.ZEVM.ETHZRC20Addr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid ETHZRC20Addr: %s", c)
		}
		r.ETHZRC20Addr = ethcommon.HexToAddress(c)
		r.ETHZRC20, err = zrc20.NewZRC20(r.ETHZRC20Addr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.ZEVM.ERC20ZRC20Addr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid ERC20ZRC20Addr: %s", c)
		}
		r.ERC20ZRC20Addr = ethcommon.HexToAddress(c)
		r.ERC20ZRC20, err = zrc20.NewZRC20(r.ERC20ZRC20Addr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.ZEVM.BTCZRC20Addr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid BTCZRC20Addr: %s", c)
		}
		r.BTCZRC20Addr = ethcommon.HexToAddress(c)
		r.BTCZRC20, err = zrc20.NewZRC20(r.BTCZRC20Addr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.ZEVM.ERC20ZRC20Addr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid ERC20ZRC20Addr: %s", c)
		}
		r.ERC20ZRC20Addr = ethcommon.HexToAddress(c)
		r.ERC20ZRC20, err = zrc20.NewZRC20(r.ERC20ZRC20Addr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.ZEVM.UniswapFactoryAddr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid UniswapFactoryAddr: %s", c)
		}
		r.UniswapV2FactoryAddr = ethcommon.HexToAddress(c)
		r.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(r.UniswapV2FactoryAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.ZEVM.UniswapRouterAddr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid UniswapRouterAddr: %s", c)
		}
		r.UniswapV2RouterAddr = ethcommon.HexToAddress(c)
		r.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(r.UniswapV2RouterAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.ZEVM.ConnectorZEVMAddr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid ConnectorZEVMAddr: %s", c)
		}
		r.ConnectorZEVMAddr = ethcommon.HexToAddress(c)
		r.ConnectorZEVM, err = connectorzevm.NewZetaConnectorZEVM(r.ConnectorZEVMAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.ZEVM.WZetaAddr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid WZetaAddr: %s", c)
		}
		r.WZetaAddr = ethcommon.HexToAddress(c)
		r.WZeta, err = wzeta.NewWETH9(r.WZetaAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.ZEVM.ZEVMSwapAppAddr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid ZEVMSwapAppAddr: %s", c)
		}
		r.ZEVMSwapAppAddr = ethcommon.HexToAddress(c)
		r.ZEVMSwapApp, err = zevmswap.NewZEVMSwapApp(r.ZEVMSwapAppAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.ZEVM.ContextAppAddr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid ContextAppAddr: %s", c)
		}
		r.ContextAppAddr = ethcommon.HexToAddress(c)
		r.ContextApp, err = contextapp.NewContextApp(r.ContextAppAddr, r.ZEVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.ZEVM.TestDappAddr; c != "" {
		if !ethcommon.IsHexAddress(c) {
			return fmt.Errorf("invalid TestDappAddr: %s", c)
		}
		r.TestDAppAddr = ethcommon.HexToAddress(c)
	}

	return nil
}
