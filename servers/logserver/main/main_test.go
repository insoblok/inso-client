package main

import (
	"eth-toy-client/config"
	contract "eth-toy-client/core/contracts"
	"eth-toy-client/core/httpapi"
	"eth-toy-client/core/logutil"
	toytypes "eth-toy-client/core/types"
	"eth-toy-client/kit/mockusdc"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

var ServerName = config.Servers.LogServer

func TestPing(t *testing.T) {
	res, err := ServerName.Ping()
	require.NoError(t, err, "❌ ping failed")
	require.Equal(t, 200, res.StatusCode, "❌ ping failed")

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err, "❌ failed to read response body")
	defer res.Body.Close()

	require.Equal(t, string(ServerName)+" says pong", string(body))
}

func TestServerConnectionRefused(t *testing.T) {
	//Need to make sure the server is not running
	res, err := ServerName.Ping()
	t.Log(err)
	t.Log(res)
	//require.NotNil(t, err, "❌ expected non-nil error when sever is not running "+pingURL)
	//require.Contains(t, err.Error(), "connection refused", "❌ expected connection refused error")
}

func TestContractAliasNotFound(t *testing.T) {
	serverConfig := ServerName.GetServerConfig()
	contractFooURL := serverConfig.GetServerUrl("contract/foo")
	res, err := http.Get(contractFooURL)
	require.NoError(t, err, "❌ contract/foo failed")

	_, apiError, err := httpapi.ParseAPIResponse[contract.DeployedContractMetaJSON](res)
	require.NoError(t, err, "❌ failed to parse response")
	require.NotNil(t, apiError, "❌ expected non-nil apiError")
	logutil.Infof("apiError: %v", apiError)
}

func TestServeInvalidUrl(t *testing.T) {
	serverConfig := ServerName.GetServerConfig()
	invalidURL := serverConfig.GetServerUrl("invalid")
	resp, err := http.Get(invalidURL)
	require.NoError(t, err, "❌ invalid url failed")
	require.Equal(t, 404, resp.StatusCode, "❌ expected not found")
	resBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "❌ expected to read response body")
	resBodyStr := strings.TrimSpace(string(resBody))
	require.Equal(t, "404 page not found", resBodyStr, "❌ expected 404 page not found")
}

func TestRegisterContractAddress(t *testing.T) {
	serverConfig := ServerName.GetServerConfig()
	registerURL := serverConfig.GetServerUrl("register-contract")
	payload := contract.DeployedContractMetaJSON{
		Alias:     "CounterV2",
		Address:   "0x1234567890123456789012345678901234567890",
		TxHash:    "0x1234567890123456789012345678901234567890",
		ABI:       mockusdc.MockusdcMetaData.ABI,
		Bytecode:  mockusdc.MockusdcMetaData.Bin,
		Timestamp: 1234567890,
		Owner:     "0x1234567890123456789012345678901234567890",
		Overwrite: true,
	}

	res, apiError, err := httpapi.PostWithAPIResponse[toytypes.AliasRegisterResponse](registerURL, payload)
	require.NoError(t, err, "❌ register contract failed")
	require.Nil(t, apiError, "❌ API Error expected to be nil")
	require.Equal(t, "ok", res.Status, "❌ register contract failed")
	require.Equal(t, "CounterV2", res.Alias, "❌ register contract failed")

	contractURL := serverConfig.GetServerUrl("contract/" + payload.Alias)
	res2, apiError, err := httpapi.PostWithAPIResponseNoPayload[contract.DeployedContractMetaJSON](contractURL)
	require.NoError(t, err, "❌ failed to get contract, expected nil error")
	require.Nil(t, apiError, "❌ failed to get contract, expected nil apiError")
	require.Equal(t, &payload, res2, "❌ payload mismatch")

}
