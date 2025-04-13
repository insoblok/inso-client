package main

import (
	"eth-toy-client/config"
	contract "eth-toy-client/core/contracts"
	"eth-toy-client/core/httpapi"
	"eth-toy-client/core/logutil"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

var ServerName = "LogServer"

func TestPing(t *testing.T) {
	serverConfig := config.GetServerConfig(ServerName)
	pingURL := serverConfig.GetServerUrl("ping")
	res, err := http.Get(pingURL)
	require.NoError(t, err, "❌ ping failed")
	require.Equal(t, 200, res.StatusCode, "❌ ping failed")

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err, "❌ failed to read response body")
	defer res.Body.Close()

	require.Equal(t, string(body), ServerName+" says pong")
}

func TestContractAliasNotFound(t *testing.T) {
	serverConfig := config.GetServerConfig(ServerName)
	contractFooURL := serverConfig.GetServerUrl("contract/foo")
	res, err := http.Get(contractFooURL)
	require.NoError(t, err, "❌ contract/foo failed")

	_, apiError, err := httpapi.ParseAPIResponse[contract.DeployedContractMetaJSON](res)
	require.NoError(t, err, "❌ failed to parse response")
	require.NotNil(t, apiError, "❌ expected non-nil apiError")
	logutil.Infof("apiError: %v", apiError)
}

func TestServerConnectionRefused(t *testing.T) {
	//Need to make sure the server is not running
	serverConfig := config.GetServerConfig(ServerName)
	pingURL := serverConfig.GetServerUrl("ping")
	_, err := http.Get(pingURL)
	t.Log(err)
	require.NotNil(t, err, "❌ expected non-nil error when sever is not running "+pingURL)
	require.Contains(t, err.Error(), "connection refused", "❌ expected connection refused error")
}

func TestServeInvalidUrl(t *testing.T) {
	serverConfig := config.GetServerConfig(ServerName)
	invalidURL := serverConfig.GetServerUrl("invalid")
	resp, err := http.Get(invalidURL)
	require.NoError(t, err, "❌ invalid url failed")
	require.Equal(t, 404, resp.StatusCode, "❌ expected not found")
	resBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "❌ expected to read response body")
	require.Equal(t, string(resBody), "404 page not found")
}
