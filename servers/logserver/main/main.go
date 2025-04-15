package main

import (
	"eth-toy-client/servers/logserver/logserver"
	"eth-toy-client/servers/servers"
)

func main() {
	logServer := &logserver.LogServer{}
	servers.StartMicroService(logServer)
	select {}
}
