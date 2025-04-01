// File: client/playground/contract_playground_test.go
package playground

import (
	"eth-toy-client/core/contracts"
	"eth-toy-client/core/devutil"
	"eth-toy-client/core/httpapi"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListRegisteredContracts(t *testing.T) {
	url := "http://localhost:8575/api/contracts"
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	data, apiErr, err := httpapi.ParseAPIResponse[[]contract.DeployedContractMetaJSON](resp)
	require.NoError(t, err)
	require.Nil(t, apiErr)
	require.NotNil(t, data)
	require.NotEmpty(t, *data)

	t.Logf("✅ Got %d registered contracts", len(*data))
	for _, meta := range *data {
		t.Logf("🔗 %s @ %s (tx: %s, owner: %s)", meta.Alias, meta.Address, meta.TxHash, meta.Owner)
	}
}

func TestGetContractByAlias(t *testing.T) {
	urls := devutil.GetUrls()
	alias := "CounterV1"
	url := fmt.Sprintf("%s/api/contracts/%s", urls.ServerURL, alias)

	t.Logf("🌐 Target URL: %s", url)

	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, apiErr, err := httpapi.ParseAPIResponse[contract.DeployedContractMetaJSON](resp)
	require.NoError(t, err)
	require.Nil(t, apiErr)
	require.NotNil(t, data)

	t.Logf("✅ Got alias %s @ %s", data.Alias, data.Address)
	t.Logf("🧾 TX: %s", data.TxHash)
	t.Logf("👤 Owner: %s", data.Owner)
	t.Logf("🧠 ABI: %s", data.ABI[:40])
}
