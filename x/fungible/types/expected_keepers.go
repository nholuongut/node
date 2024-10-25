package types

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/x/evm/statedb"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/zetacore/common"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetSequence(ctx sdk.Context, addr sdk.AccAddress) (uint64, error)
	GetModuleAccount(ctx sdk.Context, name string) types.ModuleAccountI
	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	SetAccount(ctx sdk.Context, acc types.AccountI)
}

type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

type ObserverKeeper interface {
	GetParams(ctx sdk.Context) (params observertypes.Params)
	GetSupportedChains(ctx sdk.Context) []*common.Chain
}

type EVMKeeper interface {
	ChainID() *big.Int
	GetBlockBloomTransient(ctx sdk.Context) *big.Int
	GetLogSizeTransient(ctx sdk.Context) uint64
	WithChainID(ctx sdk.Context)
	SetBlockBloomTransient(ctx sdk.Context, bloom *big.Int)
	SetLogSizeTransient(ctx sdk.Context, logSize uint64)
	EstimateGas(c context.Context, req *evmtypes.EthCallRequest) (*evmtypes.EstimateGasResponse, error)
	ApplyMessage(
		ctx sdk.Context,
		msg core.Message,
		tracer vm.EVMLogger,
		commit bool,
	) (*evmtypes.MsgEthereumTxResponse, error)
	GetAccount(ctx sdk.Context, addr ethcommon.Address) *statedb.Account
	GetCode(ctx sdk.Context, codeHash ethcommon.Hash) []byte
	SetAccount(ctx sdk.Context, addr ethcommon.Address, account statedb.Account) error
}

type AuthorityKeeper interface {
	IsAuthorized(ctx sdk.Context, address string, policyType authoritytypes.PolicyType) bool
}
