package playground

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type InfoResponse struct {
	RPCURL        string `json:"rpcUrl"`
	AccountsCount int    `json:"accountsCount"`
}

var serverURL = "http://localhost:8575"
var infoURL = serverURL + "/info"
var accountsURL = serverURL + "/accounts"

func TestPlaygroundInfo(t *testing.T) {

	resp, err := http.Get(infoURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var info InfoResponse
	err = json.Unmarshal(body, &info)
	require.NoError(t, err)

	t.Logf("ℹ️  Test server info:")
	t.Logf("   🔗 RPC URL: %s", info.RPCURL)
	t.Logf("   👤 Accounts Count: %d", info.AccountsCount)

	require.NotEmpty(t, info.RPCURL)
	require.Greater(t, info.AccountsCount, 0)
}

func TestPlaygroundAccounts(t *testing.T) {
	resp, err := http.Get(accountsURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	var accounts []struct {
		Name       string `json:"name"`
		Address    string `json:"address"`
		PrivateKey string `json:"privateKey"`
	}
	err = json.NewDecoder(resp.Body).Decode(&accounts)
	require.NoError(t, err)

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
