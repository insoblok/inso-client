package playground

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
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
