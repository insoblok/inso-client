package ethclienttest

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
	"math/big"
	"strings"
	"testing"
)

var targetHost = "http://localhost:8565"

// For the purpose of this test, it will be ok just to return the client
func SetupEthClient(t *testing.T) *ethclient.Client {
	client, err := ethclient.Dial(targetHost)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	return client
}

func SetupRpcClient(t *testing.T) *rpc.Client {
	rpcClient, err := rpc.Dial(targetHost)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	return rpcClient
}

func FetchAccounts(t *testing.T) []common.Address {
	rpcClient := SetupRpcClient(t)
	defer rpcClient.Close()

	var raw []string
	err := rpcClient.Call(&raw, "eth_accounts")
	if err != nil {
		t.Fatalf("Failed to fetch accounts: %v", err)
	}

	if len(raw) == 0 {
		t.Fatal("Expected at least one account from dev node, got zero")
	}

	var addresses []common.Address
	for _, addr := range raw {
		addresses = append(addresses, common.HexToAddress(addr))
	}
	return addresses
}

func ToBigInt(j uint64) *big.Int {
	return big.NewInt(int64(j))
}

func GetEthClientAndAccounts(t *testing.T) ([]common.Address, *ethclient.Client) {
	accounts := FetchAccounts(t)
	client := SetupEthClient(t)
	return accounts, client
}

func WithAccount(t *testing.T, f func(t *testing.T, addr common.Address, client *ethclient.Client)) {
	t.Helper()
	accounts, client := GetEthClientAndAccounts(t)
	defer client.Close()
	for _, addr := range accounts {
		f(t, addr, client)
	}
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

func TestGetCurrentBalance(t *testing.T) {
	accounts, client := GetEthClientAndAccounts(t)
	defer client.Close()
	t.Logf("✅ Fetched %d account(s):", len(accounts))
	for i, addr := range accounts {
		balance, err := client.BalanceAt(context.Background(), addr, nil)
		require.NoError(t, err)
		t.Logf("  [%d] %s: %s", i, addr, balance)
	}
}

func TestGetCurrentBalanceAtBlock(t *testing.T) {
	accounts, client := GetEthClientAndAccounts(t)
	defer client.Close()
	t.Logf("✅ Fetched %d account(s):", len(accounts))
	height, err := client.BlockNumber(context.Background())
	require.NoError(t, err)
	t.Logf("Current block height: %d", height)
	for i, addr := range accounts {
		t.Logf("account [%d] %s", i, addr)
		for j := uint64(0); j < height; j++ {
			balance, err := client.BalanceAt(context.Background(), addr, ToBigInt(j))
			require.NoError(t, err)
			t.Logf("%sblock %d: %s", strings.Repeat(" ", 8), j, balance)
		}
	}
}

func TestGetNonce(t *testing.T) {
	WithAccount(t, func(t *testing.T, addr common.Address, client *ethclient.Client) {
		nonce, err := client.NonceAt(context.Background(), addr, nil)
		require.NoError(t, err)
		t.Logf("  account %s: pending nonce-Tx counts so far %d", addr, nonce)
	})
}

func TestGetPendingNonce(t *testing.T) {
	WithAccount(t, func(t *testing.T, addr common.Address, client *ethclient.Client) {
		nonce, err := client.PendingNonceAt(context.Background(), addr)
		require.NoError(t, err)
		t.Logf("  account %s: pending nonce-Tx counts so far %d", addr, nonce)
	})
}

func TestSendTransactionEIP1559(t *testing.T) {
	WithAccount(t, func(t *testing.T, addr common.Address, client *ethclient.Client) {
		ctx := context.Background()

		// 1. Get pending nonce
		nonce, err := client.PendingNonceAt(ctx, addr)
		require.NoError(t, err)

		// 2. Set gas limit and EIP-1559 fees
		gasLimit := uint64(21000) // Standard ETH transfer

		tipCap, err := client.SuggestGasTipCap(ctx)
		require.NoError(t, err)

		// Optional: Add a little extra to ensure it's mined
		feeCap := new(big.Int).Add(tipCap, big.NewInt(2e9)) // tip + 2 gwei

		// 3. Amount to send
		amount := new(big.Int).Mul(big.NewInt(1), big.NewInt(1e18)) // 1 ETH

		// 4. Get chain ID (dev mode = 1337 by default)
		chainID, err := client.ChainID(ctx)
		require.NoError(t, err)

		// 5. Construct EIP-1559 tx
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			GasTipCap: tipCap,
			GasFeeCap: feeCap,
			Gas:       gasLimit,
			To:        &addr, // self-transfer
			Value:     amount,
			Data:      nil,
		})

		// 6. Send directly — dev mode auto-signs unlocked accounts
		err = client.SendTransaction(ctx, tx)
		require.NoError(t, err)

		t.Logf("📤 Sent tx: %s", tx.Hash().Hex())
	})
}
