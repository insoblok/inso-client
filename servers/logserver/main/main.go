package main

import (
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

func (logServer *LogServer) InitService(nodeClient *servers.NodeClient, serverConfig servers.ServerConfig) (servers.ServerConfig, http.Handler) {
	handlers := logserver.SetupRoutes()
	return serverConfig, handlers
}
