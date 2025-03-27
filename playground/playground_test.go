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
	Name       string `json:"name"`
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
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
		accountsMap[acc.Name] = acc
	}
	return accountsMap
}

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
	accounts := GetAccounts(t, GetUrls())
	require.Len(t, accounts, 10, "Expected 10 test accounts")

	alice, ok := accounts["alice"]
	require.True(t, ok, "Alice account is not found")
	bob, ok := accounts["bob"]
	require.True(t, ok, "Bob account is not found")

	resp := GetInfoResponse(t, GetUrls())

	client, err := ethclient.Dial(resp.RPCURL)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	require.NoError(t, err)
	aliceAddr := common.HexToAddress(alice.Address)
	nonce, err := client.PendingNonceAt(ctx, aliceAddr)
	require.NoError(t, err)
	bobAddr := common.HexToAddress(bob.Address)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: big.NewInt(1e9), // Max fee
		GasTipCap: big.NewInt(1),   // Priority tip
		Gas:       21_000,
		To:        &bobAddr,
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
