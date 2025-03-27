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

func GetAccounts(t *testing.T, urls Urls) []ClientTestAccount {
	resp, err := http.Get(urls.AccountsURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	var accounts []ClientTestAccount

	err = json.NewDecoder(resp.Body).Decode(&accounts)
	require.NoError(t, err)
	return accounts
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

	foundAlice := false
	for _, acc := range accounts {
		require.NotEmpty(t, acc.PrivateKey, "Account has empty private key")
		if acc.Name == "alice" {
			foundAlice = true
			t.Logf("👑 Found Alice: %s", acc.Address)
		}
	}
	require.True(t, foundAlice, "Expected to find 'alice' in test accounts")
	t.Logf("✅ Successfully fetched and verified %d accounts", len(accounts))
}
