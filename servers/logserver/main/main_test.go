package main

import (
	"eth-toy-client/config"
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
