package consts

import "math/big"
import "github.com/ethereum/go-ethereum/common"

// EthUnit holds common denominations of ETH
type EthUnit struct {
	Wei      *big.Int
	Gwei     *big.Int
	Finney   *big.Int
	Ether    *big.Int
	TenEth   *big.Int
	Point01  *big.Int
	Point001 *big.Int
}

// ETH is a globally accessible set of common units
var ETH = EthUnit{
	Wei:      big.NewInt(1),
	Gwei:     big.NewInt(1e9),
	Finney:   new(big.Int).Mul(big.NewInt(1e15), big.NewInt(1)),
	Ether:    new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1)),
	TenEth:   new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18)),
	Point01:  new(big.Int).Mul(big.NewInt(1), big.NewInt(1e16)), // 0.01 ETH
	Point001: new(big.Int).Mul(big.NewInt(1), big.NewInt(1e15)), // 0.001 ETH
}

const (
	// ⛽ Gas parameters
	GasLimitTransfer  uint64 = 21_000
	GasTipCapLow             = 1             // 1 wei
	GasFeeCapStandard        = 1_000_000_000 // 1 gwei

	// 💸 ETH Transfer amounts (in wei)
	EthAmount_001         = 1e16 // 0.01 ETH
	EthAmount_01          = 1e17 // 0.1 ETH
	EthAmount_1           = 1e18 // 1 ETH
	DefaultTransferAmount = EthAmount_001
	DefaultChainID        = 1337
)

type GasParams struct {
	GasLimitTransfer  uint64
	GasTipCapLow      uint64
	GasFeeCapStandard uint64
}

type EthAmounts struct {
	Point001           uint64 // 0.01 ETH
	Point01            uint64 // 0.1 ETH
	One                uint64 // 1 ETH
	DefaultTransferWei uint64
	DefaultChainID     int64
}

// Singleton instances
var Gas = &GasParams{
	GasLimitTransfer:  21_000,
	GasTipCapLow:      1,             // 1 wei
	GasFeeCapStandard: 1_000_000_000, // 1 gwei
}

var ETH2 = &EthAmounts{
	Point001:           1e16, // 0.01 ETH
	Point01:            1e17, // 0.1 ETH
	One:                1e18, // 1 ETH
	DefaultTransferWei: 1e16,
	DefaultChainID:     1337,
}

type CanonicalValues struct {
	ZeroAddress common.Address
	ZeroHash    common.Hash
}

var Canonical = CanonicalValues{
	ZeroAddress: common.Address{},
	ZeroHash:    common.Hash{},
}
