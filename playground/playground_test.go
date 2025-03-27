package playground

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"
)

type InfoResponse struct {
	RPCURL        string `json:"rpcUrl"`
	AccountsCount int    `json:"accountsCount"`
}

type Urls struct {
	ServerURL   string
	InfoURL     string
	AccountsURL string
}

type ClientTestAccount struct {
	Name          string         `json:"name"`
	Address       string         `json:"address"`
	PrivateKey    string         `json:"privateKey"`
	CommonAddress common.Address `json:"-"`
}

func GetUrls() Urls {
	base := "http://localhost:8575"
	return Urls{
		ServerURL:   base,
		InfoURL:     base + "/info",
		AccountsURL: base + "/accounts",
	}
}

func GetInfoResponse(t *testing.T, urls Urls) InfoResponse {
	resp, err := http.Get(urls.InfoURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var info InfoResponse
	err = json.Unmarshal(body, &info)
	require.NoError(t, err)
	return info
}

func GetAccounts(t *testing.T, urls Urls) map[string]ClientTestAccount {
	resp, err := http.Get(urls.AccountsURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	var accounts []ClientTestAccount

	err = json.NewDecoder(resp.Body).Decode(&accounts)
	require.NoError(t, err)

	accountsMap := make(map[string]ClientTestAccount)
	for _, acc := range accounts {
		acc.CommonAddress = common.HexToAddress(acc.Address)
		accountsMap[acc.Name] = acc
	}
	return accountsMap
}

func MustGet(t *testing.T, urls Urls) (*ethclient.Client, ClientTestAccount, ClientTestAccount, map[string]ClientTestAccount) {
	accounts := GetAccounts(t, urls)
	require.Len(t, accounts, 10, "Expected 10 test accounts")

	alice, ok := accounts["alice"]
	require.True(t, ok, "Alice account is not found")
	bob, ok := accounts["bob"]
	require.True(t, ok, "Bob account is not found")
	resp := GetInfoResponse(t, urls)
	client, err := ethclient.Dial(resp.RPCURL)
	require.NoError(t, err)

	return client, alice, bob, accounts

}

//////////////////////////////////////////////////////////////////

func TestPlaygroundInfo(t *testing.T) {
	info := GetInfoResponse(t, GetUrls())

	t.Logf("ℹ️  Test server info:")
	t.Logf("   🔗 RPC URL: %s", info.RPCURL)
	t.Logf("   👤 Accounts Count: %d", info.AccountsCount)
	require.NotEmpty(t, info.RPCURL)
	require.Greater(t, info.AccountsCount, 0)
}

func TestPlaygroundAccounts(t *testing.T) {

	accounts := GetAccounts(t, GetUrls())
	require.Len(t, accounts, 10, "Expected 10 test accounts")

	alice, ok := accounts["alice"]
	require.True(t, ok, "Alice account is not found")
	require.NotEmpty(t, alice.Address, "Alice's address is empty")
	require.NotEmpty(t, alice.PrivateKey, "Alice's private key is empty")
	t.Logf("🎉 Extracted Alice: Address: %s, PrivateKey: %s", alice.Address, alice.PrivateKey)

}

func TestSignedTxFromAliceToBob(t *testing.T) {
	client, alice, bob, _ := MustGet(t, GetUrls())
	defer client.Close()

	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	require.NoError(t, err)
	nonce, err := client.PendingNonceAt(ctx, alice.CommonAddress)
	require.NoError(t, err)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: big.NewInt(1e9), // Max fee
		GasTipCap: big.NewInt(1),   // Priority tip
		Gas:       21_000,
		To:        &bob.CommonAddress,
		Value:     big.NewInt(1e16), // 0.01 ETH
	})
	signer := types.NewLondonSigner(chainID)
	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(alice.PrivateKey, "0x"))
	if err != nil {
		log.Fatalf("❌ Invalid private key for %s: %v", alice.Name, err)
	}

	signedTx, err := types.SignTx(tx, signer, privKey)
	require.NoError(t, err)
	t.Logf("📤 Sent 0.01 ETH from alice to bob — tx: %s", signedTx.Hash().Hex())
}

func TestSendSignedTxFromAliceToBob(t *testing.T) {
	client, alice, bob, _ := MustGet(t, GetUrls())
	defer client.Close()

	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	require.NoError(t, err)
	nonce, err := client.PendingNonceAt(ctx, alice.CommonAddress)
	require.NoError(t, err)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: big.NewInt(1e9), // Max fee
		GasTipCap: big.NewInt(1),   // Priority tip
		Gas:       21_000,
		To:        &bob.CommonAddress,
		Value:     big.NewInt(1e16), // 0.01 ETH
	})
	signer := types.NewLondonSigner(chainID)
	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(alice.PrivateKey, "0x"))
	if err != nil {
		log.Fatalf("❌ Invalid private key for %s: %v", alice.Name, err)
	}

	signedTx, err := types.SignTx(tx, signer, privKey)
	require.NoError(t, err)

	err = client.SendTransaction(ctx, signedTx)
	require.NoError(t, err)

	t.Logf("📤 Sent 0.01 ETH from alice to bob — tx: %s", signedTx.Hash().Hex())
}

func TestSignedTxAffectsBalances(t *testing.T) {
	client, alice, bob, _ := MustGet(t, GetUrls())
	defer client.Close()

	// 💰 Balance before
	aliceBefore, _ := client.BalanceAt(context.Background(), alice.CommonAddress, nil)
	bobBefore, _ := client.BalanceAt(context.Background(), bob.CommonAddress, nil)

	// 🧾 Nonce & Chain ID
	nonce, _ := client.PendingNonceAt(context.Background(), alice.CommonAddress)
	chainID, _ := client.ChainID(context.Background())

	// 💸 Build tx
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: big.NewInt(1e9),
		GasTipCap: big.NewInt(1),
		Gas:       21000,
		To:        &bob.CommonAddress,
		Value:     big.NewInt(1e16), // 0.01 ETH
	})

	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(alice.PrivateKey, "0x"))
	signedTx, _ := types.SignTx(tx, types.NewLondonSigner(chainID), privKey)

	// 📤 Send
	err = client.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	t.Logf("Sent tx: %s", signedTx.Hash())

	// 🕰️ Wait for mining
	time.Sleep(1 * time.Second)

	// 💰 After
	aliceAfter, _ := client.BalanceAt(context.Background(), alice.CommonAddress, nil)
	bobAfter, _ := client.BalanceAt(context.Background(), bob.CommonAddress, nil)

	t.Logf("Alice: %s -> %s", aliceBefore, aliceAfter)
	t.Logf("Bob:   %s -> %s", bobBefore, bobAfter)

	require.True(t, bobAfter.Cmp(bobBefore) > 0, "Bob should have received ETH")
	require.True(t, aliceAfter.Cmp(aliceBefore) < 0, "Alice should have less due to tx + gas")
}
