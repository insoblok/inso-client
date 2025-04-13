package main

import (
	"eth-toy-client/config"
	contract "eth-toy-client/core/contracts"
	"eth-toy-client/servers/logserver/logserver"
	"eth-toy-client/servers/servers"
	"net/http"
)

func main() {
	logServer := &LogServer{}
	servers.StartMicroService(logServer)
	select {}
}

type LogServer struct{}

func (logServer *LogServer) Name() string {
	return "LogServer"
}

func (logServer *LogServer) InitService(nodeClient *servers.NodeClient, serverConfig config.ServerConfig) (config.ServerConfig, http.Handler) {
	contractRegistry := contract.NewRegistry()
	handlers := logserver.SetupRoutes(serverConfig, contractRegistry)
	return serverConfig, handlers
}
