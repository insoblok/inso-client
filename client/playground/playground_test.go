package playground

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"eth-toy-client/core/consts"
	"eth-toy-client/core/httpapi"
	toytypes "eth-toy-client/core/types"
	"eth-toy-client/sol/counter"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"
)

type InfoResponse struct {
	RPCURL        string `json:"rpcUrl"`
	AccountsCount int    `json:"accountsCount"`
}

type Urls struct {
	ServerURL   string
	InfoURL     string
	AccountsURL string
}

type ClientTestAccount struct {
	Name          string         `json:"name"`
	Address       string         `json:"address"`
	PrivateKey    string         `json:"privateKey"`
	CommonAddress common.Address `json:"-"`
}

func GetUrls() Urls {
	base := "http://localhost:8575"
	return Urls{
		ServerURL:   base,
		InfoURL:     base + "/info",
		AccountsURL: base + "/accounts",
	}
}

func GetInfoResponse(t *testing.T, urls Urls) InfoResponse {
	resp, err := http.Get(urls.InfoURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var info InfoResponse
	err = json.Unmarshal(body, &info)
	require.NoError(t, err)
	return info
}

func GetAccounts(t *testing.T, urls Urls) map[string]ClientTestAccount {
	resp, err := http.Get(urls.AccountsURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	var accounts []ClientTestAccount

	err = json.NewDecoder(resp.Body).Decode(&accounts)
	require.NoError(t, err)

	accountsMap := make(map[string]ClientTestAccount)
	for _, acc := range accounts {
		acc.CommonAddress = common.HexToAddress(acc.Address)
		accountsMap[acc.Name] = acc
	}
	return accountsMap
}

func MustGet(t *testing.T, urls Urls) (*ethclient.Client, ClientTestAccount, ClientTestAccount, map[string]ClientTestAccount) {
	accounts := GetAccounts(t, urls)
	require.Len(t, accounts, 10, "Expected 10 test accounts")

	alice, ok := accounts["alice"]
	require.True(t, ok, "Alice account is not found")
	bob, ok := accounts["bob"]
	require.True(t, ok, "Bob account is not found")
	resp := GetInfoResponse(t, urls)
	client, err := ethclient.Dial(resp.RPCURL)
	require.NoError(t, err)

	return client, alice, bob, accounts

}

func AliceSignAndSendTx(t *testing.T) (*ethclient.Client, *types.Transaction, ClientTestAccount, ClientTestAccount, *big.Int, *big.Int) {
	client, alice, bob, _ := MustGet(t, GetUrls())

	aliceBefore, _ := client.BalanceAt(context.Background(), alice.CommonAddress, nil)
	bobBefore, _ := client.BalanceAt(context.Background(), bob.CommonAddress, nil)
	nonce, _ := client.PendingNonceAt(context.Background(), alice.CommonAddress)
	chainID, _ := client.ChainID(context.Background())
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: big.NewInt(1),
		GasFeeCap: big.NewInt(1e9),
		Gas:       21_000,
		To:        &bob.CommonAddress,
		Value:     big.NewInt(1e16), // 0.01 ETH
	})

	privKey, _ := crypto.HexToECDSA(strings.TrimPrefix(alice.PrivateKey, "0x"))
	signedTx, _ := types.SignTx(tx, types.NewLondonSigner(chainID), privKey)
	err := client.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	t.Logf("📤 Sent tx: %s", signedTx.Hash())
	return client, signedTx, alice, bob, aliceBefore, bobBefore
}

//////////////////////////////////////////////////////////////////

func TestPlaygroundInfo(t *testing.T) {
	info := GetInfoResponse(t, GetUrls())

	t.Logf("ℹ️  Test server info:")
	t.Logf("   🔗 RPC URL: %s", info.RPCURL)
	t.Logf("   👤 Accounts Count: %d", info.AccountsCount)
	require.NotEmpty(t, info.RPCURL)
	require.Greater(t, info.AccountsCount, 0)
}

func TestPlaygroundAccounts(t *testing.T) {

	accounts := GetAccounts(t, GetUrls())
	require.Len(t, accounts, 10, "Expected 10 test accounts")

	alice, ok := accounts["alice"]
	require.True(t, ok, "Alice account is not found")
	require.NotEmpty(t, alice.Address, "Alice's address is empty")
	require.NotEmpty(t, alice.PrivateKey, "Alice's private key is empty")
	t.Logf("🎉 Extracted Alice: Address: %s, PrivateKey: %s", alice.Address, alice.PrivateKey)

}

func TestSignedTxFromAliceToBob(t *testing.T) {
	client, alice, bob, _ := MustGet(t, GetUrls())
	defer client.Close()

	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	require.NoError(t, err)
	nonce, err := client.PendingNonceAt(ctx, alice.CommonAddress)
	require.NoError(t, err)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: big.NewInt(1e9), // Max fee
		GasTipCap: big.NewInt(1),   // Priority tip
		Gas:       21_000,
		To:        &bob.CommonAddress,
		Value:     big.NewInt(1e16), // 0.01 ETH
	})
	signer := types.NewLondonSigner(chainID)
	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(alice.PrivateKey, "0x"))
	if err != nil {
		log.Fatalf("❌ Invalid private key for %s: %v", alice.Name, err)
	}

	signedTx, err := types.SignTx(tx, signer, privKey)
	require.NoError(t, err)
	t.Logf("📤 Sent 0.01 ETH from alice to bob — tx: %s", signedTx.Hash().Hex())
}

func TestSendSignedTxFromAliceToBob(t *testing.T) {
	client, _, _, _, _, _ := AliceSignAndSendTx(t)
	defer client.Close()
}

func TestSignedTxAffectsBalances(t *testing.T) {
	client, _, alice, bob, aliceBefore, bobBefore := AliceSignAndSendTx(t)
	defer client.Close()

	time.Sleep(1 * time.Second)

	aliceAfter, _ := client.BalanceAt(context.Background(), alice.CommonAddress, nil)
	bobAfter, _ := client.BalanceAt(context.Background(), bob.CommonAddress, nil)

	t.Logf("Alice: %s -> %s", aliceBefore, aliceAfter)
	t.Logf("Bob:   %s -> %s", bobBefore, bobAfter)

	require.True(t, bobAfter.Cmp(bobBefore) > 0, "Bob should have received ETH")
	require.True(t, aliceAfter.Cmp(aliceBefore) < 0, "Alice should have less due to tx + gas")
}

func TestQueryTxByHash(t *testing.T) {
	client, signedTx, _, bob, _, _ := AliceSignAndSendTx(t)
	defer client.Close()

	time.Sleep(1 * time.Second)

	txBack, isPending, err := client.TransactionByHash(context.Background(), signedTx.Hash())
	require.NoError(t, err)
	require.False(t, isPending, "Transaction is still pending")

	require.Equal(t, bob.CommonAddress, *txBack.To(), "To address mismatch")
	require.Equal(t, big.NewInt(1e16), txBack.Value(), "Transferred value mismatch")
	t.Log("✅ Tx confirmed in block and matches expected recipient + value")
}

func TestTxReceiptShowsSuccess(t *testing.T) {
	client, signedTx, _, _, _, _ := AliceSignAndSendTx(t)
	defer client.Close()

	var receipt *types.Receipt
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		receipt, _ = client.TransactionReceipt(context.Background(), signedTx.Hash())
		if receipt != nil {
			break
		}
		t.Log("⏳ Waiting for receipt...")
	}
	require.NotNil(t, receipt, "Did not receive a receipt in time")
	require.Equal(t, uint64(1), receipt.Status, "Transaction failed")

	t.Logf("✅ Mined in block %d — gas used: %d", receipt.BlockNumber.Uint64(), receipt.GasUsed)
}

func TestGenerateWalletKey(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	privateKeyBytes := crypto.FromECDSA(privateKey)
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	publicKeyBytes := crypto.FromECDSAPub(publicKey)
	address := crypto.PubkeyToAddress(*publicKey)

	fmt.Printf("🔐 Private Key: 0x%x\n", privateKeyBytes)
	fmt.Printf("🔓 Public Key: 0x%x\n", publicKeyBytes)
	fmt.Printf("📮 Address: %s\n", address.Hex())
}

func TestRecoverPublicKeyFromSignature(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	message := []byte("Test message for signing")
	hash := crypto.Keccak256Hash(message)

	// Sign the message
	sig, err := crypto.Sign(hash.Bytes(), privateKey)
	require.NoError(t, err)

	// Recover public key
	recovered, err := crypto.SigToPub(hash.Bytes(), sig)
	require.NoError(t, err)

	originalBytes := crypto.FromECDSAPub(publicKey)
	recoveredBytes := crypto.FromECDSAPub(recovered)

	t.Logf("🔐 Original PubKey:  %x", originalBytes)
	t.Logf("🧠 Recovered PubKey: %x", recoveredBytes)

	require.Equal(t, originalBytes, recoveredBytes, "Recovered public key should match original")
}

func TestSignTxFromAlice(t *testing.T) {
	urls := GetUrls()
	accountMap := GetAccounts(t, urls)

	alice := accountMap["alice"]
	bob := accountMap["bob"]

	payload := map[string]any{
		"from":  alice.Name,
		"to":    bob.Address,
		"value": "10000000000000000", // 0.01 ETH
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(urls.ServerURL+"/sign-tx", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	var result struct {
		Tx string `json:"tx"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	t.Logf("🖋️ Signed Tx: %s", result.Tx)
	require.True(t, strings.HasPrefix(result.Tx, "0x"), "Expected hex-encoded tx")
	require.Greater(t, len(result.Tx), 10, "Tx should be non-trivial")
}

func TestSendTxViaDevServer(t *testing.T) {
	urls := GetUrls()
	client, alice, bob, _ := MustGet(t, urls)
	defer client.Close()

	// 💰 Get balances before
	aliceBefore, _ := client.BalanceAt(context.Background(), alice.CommonAddress, nil)
	bobBefore, _ := client.BalanceAt(context.Background(), bob.CommonAddress, nil)

	// 📨 Request send-tx
	payload := map[string]string{
		"from": alice.Name,
		"to":   bob.Name,
	}
	data, _ := json.Marshal(payload)
	resp, err := http.Post(urls.ServerURL+"/send-tx", "application/json", bytes.NewReader(data))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var txResp struct {
		TxHash string `json:"txHash"`
	}
	err = json.NewDecoder(resp.Body).Decode(&txResp)
	require.NoError(t, err)
	t.Logf("📤 Sent tx from alice to bob via dev server: %s", txResp.TxHash)

	// ⏳ Wait for mining
	time.Sleep(1 * time.Second)

	// 💰 Get balances after
	aliceAfter, _ := client.BalanceAt(context.Background(), alice.CommonAddress, nil)
	bobAfter, _ := client.BalanceAt(context.Background(), bob.CommonAddress, nil)

	t.Logf("Alice: %s → %s", aliceBefore, aliceAfter)
	t.Logf("Bob:   %s → %s", bobBefore, bobAfter)

	require.True(t, aliceAfter.Cmp(aliceBefore) < 0, "Alice should have less ETH")
	require.True(t, bobAfter.Cmp(bobBefore) > 0, "Bob should have received ETH")
}

func TestSignTxViaDevServerAPI(t *testing.T) {
	urls := GetUrls()
	client, alice, bob, _ := MustGet(t, urls)
	defer client.Close()

	req := map[string]string{
		"from":  alice.Name,
		"to":    bob.Name,
		"value": consts.ETH.Point01.String(), // 0.01 ETH
	}

	apiResp, apiErr, err := httpapi.PostWithAPIResponse[toytypes.SignTxAPIResponse](urls.ServerURL+"/api/sign-tx", req)
	require.NoError(t, err)

	require.Nil(t, apiErr)
	require.NotNil(t, apiResp.SignedTx)
	require.NotEmpty(t, apiResp.TxHash)

	t.Logf("🖋️ Signed tx from API: hash=%s", apiResp.TxHash)
}

func TestSendTxViaDevServerAPI(t *testing.T) {
	urls := GetUrls()
	client, alice, bob, _ := MustGet(t, urls)
	defer client.Close()

	req := toytypes.SignTxRequest{
		From:  alice.Name,
		To:    bob.Name,
		Value: consts.ETH.Point01.String(),
	}

	apiResp, apiErr, err := httpapi.PostWithAPIResponse[toytypes.SendTxAPIResponse](urls.ServerURL+"/api/send-tx", req)
	require.NoError(t, err)

	if apiErr != nil {
		t.Fatalf("❌ API Error: %s — %s", apiErr.Code, apiErr.Message)
	}

	require.NotEmpty(t, apiResp.TxHash)
	t.Logf("📤 Sent tx from Alice to Bob via API: %s", apiResp.TxHash)

	// ⛏ Confirm tx impact
	time.Sleep(1 * time.Second)
	aliceAfter, _ := client.BalanceAt(context.Background(), alice.CommonAddress, nil)
	bobAfter, _ := client.BalanceAt(context.Background(), bob.CommonAddress, nil)

	t.Logf("Alice: %s", aliceAfter)
	t.Logf("Bob:   %s", bobAfter)
}

func TestDeployCounterContractViaAPI(t *testing.T) {
	urls := GetUrls()

	// ✅ Use our MustGet helper for unified client and accounts
	client, alice, _, _ := MustGet(t, urls)
	defer client.Close()

	// 🧱 Load contract bytecode from Go binding (CounterMetaData)
	data := strings.TrimPrefix(counter.CounterMetaData.Bin, "0x")
	bytecode := strings.TrimSpace(data)
	if !strings.HasPrefix(bytecode, "0x") {
		bytecode = "0x" + bytecode
	}

	t.Logf("📦 Deploying contract using bytecode (length %d bytes)", len(bytecode)/2)

	// 📨 Compose SignTxRequest (no "To" field means contract deployment)
	req := toytypes.SignTxRequest{
		From:  alice.Name,
		To:    "", // No recipient for contract deployment
		Value: "0",
		Data:  bytecode,
	}

	// 📤 Send the request to DevServer
	apiResp, apiErr, err := httpapi.PostWithAPIResponse[toytypes.SendTxAPIResponse](urls.ServerURL+"/api/send-tx", req)
	require.NoError(t, err)
	require.Nil(t, apiErr)
	require.NotEmpty(t, apiResp.TxHash)

	t.Logf("🚀 Contract deployment tx hash: %s", apiResp.TxHash)

	// ⏳ Wait for transaction to be mined and get the receipt
	receipt := WaitForReceipt(t, client, common.HexToHash(apiResp.TxHash))
	require.NotNil(t, receipt)

	t.Logf("✅ Contract deployed at: %s (block %d)", receipt.ContractAddress.Hex(), receipt.BlockNumber.Uint64())
}

func WaitForReceipt(t *testing.T, client *ethclient.Client, txHash common.Hash) *types.Receipt {
	ctx := context.Background()
	for i := 0; i < 60; i++ {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			t.Logf("🧾 Receipt received after %d sec", i)
			return receipt
		}
		t.Logf("⏳ Waiting for receipt... (%d)", i)
		time.Sleep(1 * time.Second)
	}
	t.Fatalf("⏱️ Timeout waiting for receipt of tx %s", txHash.Hex())
	return nil
}
