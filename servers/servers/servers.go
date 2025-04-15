package servers

import (
	"eth-toy-client/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"math/big"
	"net/http"
	"time"
)

type ChainId *big.Int
type Nonce *big.Int

func PingDevNode(rpcClient *rpc.Client) bool {
	var result string
	err := rpcClient.Call(&result, "web3_clientVersion")
	return err == nil
}

func ConnectToDevNode(config config.DevNodeConfig) (*rpc.Client, <-chan struct{}, error) {
	client, err := rpc.Dial("http://localhost:" + config.Port)
	if err != nil {
		return nil, nil, err
	}

	ready := make(chan struct{})
	go func() {
		for {
			if PingDevNode(client) {
				log.Printf("✅ Geth dev node is ready on port %s", config.Port)
				close(ready)
				return
			}
			log.Println("⏳ Waiting for Geth to be ready...")
			time.Sleep(1 * time.Second)
		}
	}()

	return client, ready, nil
}

type NodeClient struct {
	Config    config.DevNodeConfig
	Client    *ethclient.Client
	RPCClient *rpc.Client
	WSClient  *ethclient.Client
}

func EstablishConnectionToDevNode(name config.ServerName) (config.ServerConfig, *NodeClient) {
	serverConfig := name.GetServerConfig()
	log.Printf("📡 starting Server: %+v", serverConfig)
	rpcClient, readyChannel, err := ConnectToDevNode(serverConfig.DevNodeConfig)
	if err != nil {
		log.Fatalf("Error starting dev node: %v", err)
	}

	select {
	case <-readyChannel:
		log.Println("🚦 Node is readyChannel. Proceed.")
	case <-time.After(5 * time.Second):
		log.Fatal("🕒 Timeout waiting for dev node to start.")
	}

	client := ethclient.NewClient(rpcClient)
	wsClient, err := ethclient.Dial("ws://127.0.0.1:" + serverConfig.DevNodeConfig.WebSocketPort)
	if err != nil {
		log.Fatalf("❌ Failed to connect to WebSocket: %v", err)
	}
	log.Println("✅ Connected to Geth via WebSocket")

	return serverConfig,
		&NodeClient{
			Config:    serverConfig.DevNodeConfig,
			Client:    client,
			RPCClient: rpcClient,
			WSClient:  wsClient,
		}
}

type MicroService interface {
	Name() config.ServerName
	InitService(nodeClient *NodeClient, serverConfig config.ServerConfig) (config.ServerConfig, http.Handler)
}

func StartMicroService(microService MicroService) {
	serverConfig, nodeClient := EstablishConnectionToDevNode(microService.Name())
	defer nodeClient.Client.Close()
	defer nodeClient.RPCClient.Close()
	_, handler := microService.InitService(nodeClient, serverConfig)
	go func() {
		log.Println("🌐 " + string(serverConfig.Name) + " You can ping http://localhost:" + serverConfig.Port + "/ping ...")
		err := http.ListenAndServe(":"+serverConfig.Port, handler)
		if err != nil {
			log.Fatalf("❌ Failed to start HTTP server: %v", err)
		}
	}()
}
