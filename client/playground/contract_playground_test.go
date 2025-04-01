// File: client/playground/contract_playground_test.go
package playground

import (
	"eth-toy-client/core/contracts"
	"eth-toy-client/core/httpapi"
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
