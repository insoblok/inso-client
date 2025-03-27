package devnode

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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
