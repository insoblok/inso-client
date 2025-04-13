package servers

import (
	"flag"
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

type DevNodeConfig struct {
	Port string
}

type ServerConfig struct {
	Name          string
	Port          string
	DevNodeConfig DevNodeConfig
}

func GetServerConfig(name string) ServerConfig {
	var port string
	var serverPort string
	flag.StringVar(&port, "port", "8545", "HTTP RPC port for the dev node")
	flag.StringVar(&serverPort, "serverPort", "8888", "HTTP RPC port for the supporting server")
	flag.Parse()

	return ServerConfig{
		Name: name,
		Port: serverPort,
		DevNodeConfig: DevNodeConfig{
			Port: port,
		},
	}
}

func ConnectToDevNode(config DevNodeConfig) (*rpc.Client, <-chan struct{}, error) {
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
	Config    DevNodeConfig
	Client    *ethclient.Client
	RPCClient *rpc.Client
}

func EstablishConnectionToDevNode() (ServerConfig, *NodeClient) {
	serverConfig := GetServerConfig("DevServer")
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
	return serverConfig,
		&NodeClient{
			Config:    serverConfig.DevNodeConfig,
			Client:    client,
			RPCClient: rpcClient,
		}
}

type MicroService interface {
	InitService(nodeClient *NodeClient, serverConfig ServerConfig) (ServerConfig, http.Handler)
}

func StartMicroService(microService MicroService) {
	serverConfig, nodeClient := EstablishConnectionToDevNode()
	defer nodeClient.Client.Close()
	defer nodeClient.RPCClient.Close()
	_, handler := microService.InitService(nodeClient, serverConfig)
	go func() {
		log.Println("🌐 " + serverConfig.Name + "listening at http://localhost:" + serverConfig.Port + "...")
		err := http.ListenAndServe(":"+serverConfig.Port, handler)
		if err != nil {
			log.Fatalf("❌ Failed to start HTTP server: %v", err)
		}
	}()
}
