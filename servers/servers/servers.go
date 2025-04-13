package servers

import (
	"flag"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"math/big"
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
	RPCPort string
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
			RPCPort: port,
		},
	}
}

func ConnectToDevNode(config DevNodeConfig) (*rpc.Client, <-chan struct{}, error) {
	client, err := rpc.Dial("http://localhost:" + config.RPCPort)
	if err != nil {
		return nil, nil, err
	}

	ready := make(chan struct{})
	go func() {
		for {
			if PingDevNode(client) {
				log.Printf("✅ Geth dev node is ready on port %s", config.RPCPort)
				close(ready)
				return
			}
			log.Println("⏳ Waiting for Geth to be ready...")
			time.Sleep(1 * time.Second)
		}
	}()

	return client, ready, nil
}

func EstablishConnectionToDevNode() (ServerConfig, *ethclient.Client) {
	serverConfig := GetServerConfig("DevServer")
	log.Printf("starting Server: %+v", serverConfig)
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
	return serverConfig, client
}
