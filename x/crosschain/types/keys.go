package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName defines the module name
	ModuleName = "crosschain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_metacore"

	ProtocolFee = 2000000000000000000
	//TssMigrationGasMultiplierEVM is multiplied to the median gas price to get the gas price for the tss migration . This is done to avoid the tss migration tx getting stuck in the mempool
	TssMigrationGasMultiplierEVM = "2.5"
)

func GetProtocolFee() sdk.Uint {
	return sdk.NewUint(ProtocolFee)
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	SendKey              = "Send-value-"
	LastBlockHeightKey   = "LastBlockHeight-value-"
	FinalizedInboundsKey = "FinalizedInbounds-value-"

	GasPriceKey = "GasPrice-value-"

	GasBalanceKey = "GasBalance-value-"

	OutTxTrackerKeyPrefix = "OutTxTracker-value-"
	InTxTrackerKeyPrefix  = "InTxTracker-value-"

	// #nosec G101: Potential hardcoded credentials (gosec)
	// ZetaAccountingKey value is used as prefix for storing ZetaAccountingKey
	ZetaAccountingKey = "ZetaAccounting-value-"
)

// OutTxTrackerKey returns the store key to retrieve a OutTxTracker from the index fields
func OutTxTrackerKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}

// TODO: what's the purpose of this log identifier?
func (m CrossChainTx) LogIdentifierForCCTX() string {
	if len(m.OutboundTxParams) == 0 {
		return fmt.Sprintf("%s-%d", m.InboundTxParams.Sender, m.InboundTxParams.SenderChainId)
	}
	i := len(m.OutboundTxParams) - 1
	outTx := m.OutboundTxParams[i]
	return fmt.Sprintf("%s-%d-%d-%d", m.InboundTxParams.Sender, m.InboundTxParams.SenderChainId, outTx.ReceiverChainId, outTx.OutboundTxTssNonce)

}

func FinalizedInboundKey(intxHash string, chainID int64, eventIndex uint64) string {
	return fmt.Sprintf("%d-%s-%d", chainID, intxHash, eventIndex)
}

var (
	ModuleAddress = authtypes.NewModuleAddress(ModuleName)
	//ModuleAddressEVM common.EVMAddress
	ModuleAddressEVM = common.BytesToAddress(ModuleAddress.Bytes())
	//0xB73C0Aac4C1E606C6E495d848196355e6CB30381
)
