package servers

import (
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
