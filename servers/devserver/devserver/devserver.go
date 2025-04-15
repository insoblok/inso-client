package devserver

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"eth-toy-client/servers/servers"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"math/big"
	"strings"
)

type TestAccount struct {
	Address common.Address
	Name    string
	PrivKey *ecdsa.PrivateKey
}

var rawAccounts = map[string]string{
	"alice":   "0x4f3edf983ac636a65a842ce7c78d9aa706d3b113b37b1e1a7e3975d3fef1fdc3",
	"bob":     "0x6c8750b42f08aabb2645f65f38d013ce59b2a2b3d2a2b9c54b541714aa30353b",
	"charlie": "0x5c7d2381d08d6579e9af64a2b7d00e33a1f519896f17f0bb3e2ff894f9b8c67f",
	"diana":   "0x646f1ce2b44b655e5d34f4313ffb4c7b06d93b219a22993b7824b8b36aa5488e",
	"eric":    "0xadd53f9a7e588d21953d7cc5c5e1e89f9dfd4e0f8ee0b0b98f71f5b38d9e3a3e",
	"frank":   "0x395df67f0c2f20f3c3f4817c88c7b4c2dbd0493a34f52fd0e9f6f0017cba94d7",
	"grace":   "0xe485d098507f9e4cdbf3d0107c4b74757b6a12a7b1300a23a1f2fffae6d1c0a7",
	"helen":   "0xa453611d9419d0e3fdae7d16b1b5c470fe15f0ff9891753c603a00f226b84aa1",
	"ivan":    "0x69ee0de5eae9c6f3b674e41c2c425b6f1d50dfe2e6ec8ab3403d0fd9e3f84b07",
	"judy":    "0xdbda1821b80551c9e1c03fa0c27f14f85f3bdd8e991b52e3b0a38a42f4c1d2a3",
}

func LoadTestAccounts() *map[string]*TestAccount {
	var testAccounts = make(map[string]*TestAccount)
	for name, hexKey := range rawAccounts {
		privKey, err := crypto.HexToECDSA(strings.TrimPrefix(hexKey, "0x"))
		if err != nil {
			log.Fatalf("❌ Invalid private key for %s: %v", name, err)
		}
		addr := crypto.PubkeyToAddress(privKey.PublicKey)
		testAccounts[name] = &TestAccount{
			Name:    name,
			Address: addr,
			PrivKey: privKey,
		}
	}

	log.Printf("🧾 Loaded %d test accounts:", len(testAccounts))
	for name, acc := range testAccounts {
		log.Printf("  %s => %s", name, acc.Address.Hex())
	}
	return &testAccounts
}

func FundTestAccounts(
	devAccount common.Address,
	rpcClient *rpc.Client,
	testAccounts *map[string]*TestAccount) *map[string]*TestAccount {
	ctx := context.Background()

	for name, acc := range *testAccounts {
		// Build a call to eth_sendTransaction with from, to, value (in hex)
		var txHash string
		err := rpcClient.CallContext(ctx, &txHash, "eth_sendTransaction", map[string]interface{}{
			"from":  devAccount.Hex(),
			"to":    acc.Address.Hex(),
			"value": "0xde0b6b3a7640000", // 1 ETH = 10^18 = 0xde0b6b3a7640000
		})
		if err != nil {
			log.Fatalf("❌ eth_sendTransaction failed for %s: %v", name, err)
		}

		log.Printf("📤 Funded %s (%s) with 1 ETH (tx: %s)", name, acc.Address.Hex(), txHash)
	}
	return testAccounts
}

type SignTxResponse struct {
	Tx string `json:"tx"` // signed RLP hex
}

type SendTxResponse struct {
	TxHash string `json:"txHash"`
}

func BuildAndSignTx(
	privKey *ecdsa.PrivateKey,
	from common.Address,
	to *common.Address, // ✅ nil means contract deployment
	value *big.Int,
	rpcPort string,
	data []byte, // ✅ Optional data (contract bytecode or calldata)
) (*types.Transaction, *types.Transaction, error) {
	client, err := ethclient.Dial("http://localhost:" + rpcPort)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to dev node: %w", err)
	}
	defer client.Close()

	ctx := context.Background()

	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	address := crypto.CreateAddress(from, nonce)
	log.Printf(
		"Expected address: %s\n", address)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: big.NewInt(1),
		GasFeeCap: big.NewInt(1_000_000_000), // 1 gwei
		Gas:       3_000_000,                 // ⛽ for deployment or interaction
		To:        to,
		Value:     big.NewInt(0),
		Data:      data, // 🧠 smart contract bytecode or calldata
	})

	signedTx, err := SignTx(chainID, tx, privKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign tx: %w", err)
	}

	return tx, signedTx, nil
}

func SignContract(
	privKey *ecdsa.PrivateKey,
	from common.Address,
	rpcPort string,
	data []byte, // ✅ Optional data (contract bytecode or calldata)
) (*types.Transaction, *common.Address, *types.Transaction, error) {
	client, err := ethclient.Dial("http://localhost:" + rpcPort)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to dev node: %w", err)
	}
	defer client.Close()

	ctx := context.Background()

	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: big.NewInt(1),
		GasFeeCap: big.NewInt(1_000_000_000), // 1 gwei
		Gas:       3_000_000,                 // ⛽ for deployment or interaction
		To:        nil,
		Value:     big.NewInt(0),
		Data:      data, // 🧠 smart contract bytecode or calldata
	})

	signedTx, err := SignTx(chainID, tx, privKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to sign tx: %w", err)
	}

	address := crypto.CreateAddress(from, nonce)
	log.Printf(
		"Expected address: %s\n", address)
	return tx, &address, signedTx, nil
}

func SignTx(chainID servers.ChainId, rawTx *types.Transaction, key *ecdsa.PrivateKey) (*types.Transaction, error) {
	singer := types.NewPragueSigner(chainID)
	digest := singer.Hash(rawTx).Bytes()
	fmt.Printf("🦄 digestLength: %x\n", len(digest))
	signature, err := crypto.Sign(digest, key)
	if err != nil {
		return nil, err
	}
	return rawTx.WithSignature(
		singer,
		signature)

}

// RlpEncodeBytes returns raw RLP-encoded tx bytes
func RlpEncodeBytes(tx *types.Transaction) []byte {
	var buf bytes.Buffer
	if err := rlp.Encode(&buf, tx); err != nil {
		log.Fatalf("❌ Failed to RLP-encode tx: %v", err)
	}
	return buf.Bytes()
}

// RlpEncodeHex returns the RLP-encoded tx as hex string (0x-prefixed)
func RlpEncodeHex(tx *types.Transaction) string {
	return "0x" + hex.EncodeToString(RlpEncodeBytes(tx))
}
