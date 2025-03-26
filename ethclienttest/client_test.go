package ethclienttest

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"testing"
)

// For the purpose of this test, it will be ok just to return the client
func SetupEthClient(t *testing.T) *ethclient.Client {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	return client
}

func SetupRpcClient(t *testing.T) *rpc.Client {
	rpcClient, err := rpc.Dial("http://localhost:8545")
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	return rpcClient
}

func FetchAccounts(t *testing.T) []string {
	rpcClient := SetupRpcClient(t)
	defer rpcClient.Close()

	var accounts []string
	err := rpcClient.Call(&accounts, "eth_accounts")
	if err != nil {
		t.Fatalf("Failed to fetch accounts: %v", err)
	}

	if len(accounts) == 0 {
		t.Fatal("Expected at least one account from dev node, got zero")
	}
	return accounts
}

/////////////// Test func ///////////////

func TestDailFailWithUnknownScheme(t *testing.T) {
	client, err := ethclient.Dial("foo://localhost:8545")
	if client != nil {
		t.Fatalf("We have a client ????: %v", client)
	}
	t.Logf("We failed to dial %s", err)
}

func TestConnectToNode(t *testing.T) {
	client := SetupEthClient(t)
	defer client.Close()
	t.Log("Connected to node")
}

func TestBlockNumber(t *testing.T) {
	client := SetupEthClient(t)
	defer client.Close()
	num, err := client.BlockNumber(context.Background())
	if err != nil {
		t.Fatalf("Failed to get block number: %v", err)
	}

	t.Logf("Block number: %d", num)
}

func TestGetAccountsFromNode(t *testing.T) {
	accounts := FetchAccounts(t)

	t.Logf("✅ Fetched %d account(s):", len(accounts))
	for i, addr := range accounts {
		t.Logf("  [%d] %s", i, addr)
	}
}

func TestGetBalance(t *testing.T) {
	accounts := FetchAccounts(t)
	client := SetupEthClient(t)
	defer client.Close()
	t.Logf("✅ Fetched %d account(s):", len(accounts))
	for i, addr := range accounts {
		client.BalanceAt(context.Background(), addr, nil)
	}

}

//func TestBlockHeightIncreasesAfterTx(t *testing.T) {
//	client := SetupEthClient(t)
//	defer client.Close()
//	currentHeight, err := client.BlockNumber(context.Background())
//	require.NoError(t, err)
//	t.Logf("Current block height: %d", currentHeight)
//	tx := ???
//	client.SendTransaction(context.Background(), nil)
//
//}
