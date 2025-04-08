package accounts

import (
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"testing"
)

func TestGenerateAccount(t *testing.T) {

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("❌ Failed to generate key: %v", err)
	}
	t.Logf("✅ keyGnerated: %v", key)
	t.Logf("✅ public: %v", key.PublicKey)
	t.Logf("✅ D: %v", key.D)
	t.Logf("✅ X: %v", key.X)
	t.Logf("✅ Y: %v", key.Y)
	t.Logf("✅ Address: %v", crypto.PubkeyToAddress(key.PublicKey).Hex())
}

func TestBigInt(t *testing.T) {
	a := big.NewInt(0)
	a.SetString("237529037524875297532557752895265289659659", 10)
	x := a.Text(16)
	t.Logf("hex: 0x%s", x)
	key, err := crypto.HexToECDSA(x)
	if key != nil {
		log.Fatalf("❌ How could this happend %v", key)
	}
	log.Printf("err: %v", err)

}

func TestHexToPrivateKey(t *testing.T) {
	x := "0xadd53f9a7e588d21953d7cc5c5e1e89f9dfd4e0f8ee0b0b98f71f5b38d9e3a3e"
	xTrimmed := x[2:]
	log.Printf("Len: %v", len(x))
	log.Printf("Len Trimmed: %v", xTrimmed)
	key, err := crypto.HexToECDSA(xTrimmed)
	if err != nil {
		log.Fatalf("❌ Failed to convert hex to key: %v", err)
	}
	log.Printf("key: %v", key)
	address := crypto.PubkeyToAddress(key.PublicKey)
	log.Printf("✅ address: %v", address)
}
