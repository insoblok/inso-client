package playground

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"eth-toy-client/core/consts"
	contract "eth-toy-client/core/contracts"
	"eth-toy-client/core/devutil"
	"eth-toy-client/core/httpapi"
	"eth-toy-client/core/logutil"
	toytypes "eth-toy-client/core/types"
	"eth-toy-client/sol/out/counter"
	"fmt"
	"github.com/ethereum/go-ethereum"
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

func GetInfoResponse(t *testing.T, urls devutil.Urls) devutil.InfoResponse {
	resp, err := http.Get(urls.InfoURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var info devutil.InfoResponse
	err = json.Unmarshal(body, &info)
	require.NoError(t, err)
	return info
}

func GetAccounts(t *testing.T, urls devutil.Urls) map[string]devutil.ClientTestAccount {
	resp, err := http.Get(urls.AccountsURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	var accounts []devutil.ClientTestAccount

	err = json.NewDecoder(resp.Body).Decode(&accounts)
	require.NoError(t, err)

	accountsMap := make(map[string]devutil.ClientTestAccount)
	for _, acc := range accounts {
		acc.CommonAddress = common.HexToAddress(acc.Address)
		accountsMap[acc.Name] = acc
	}
	return accountsMap
}

func MustGet(t *testing.T, urls devutil.Urls) (*ethclient.Client, devutil.ClientTestAccount, devutil.ClientTestAccount, map[string]devutil.ClientTestAccount) {
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

func AliceSignAndSendTx(t *testing.T) (*ethclient.Client, *types.Transaction, devutil.ClientTestAccount, devutil.ClientTestAccount, *big.Int, *big.Int) {
	client, alice, bob, _ := MustGet(t, devutil.GetUrls())

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
	info := GetInfoResponse(t, devutil.GetUrls())

	t.Logf("ℹ️  Test server info:")
	t.Logf("   🔗 RPC URL: %s", info.RPCURL)
	t.Logf("   👤 Accounts Count: %d", info.AccountsCount)
	require.NotEmpty(t, info.RPCURL)
	require.Greater(t, info.AccountsCount, 0)
}

func TestPlaygroundAccounts(t *testing.T) {

	accounts := GetAccounts(t, devutil.GetUrls())
	require.Len(t, accounts, 10, "Expected 10 test accounts")

	alice, ok := accounts["alice"]
	require.True(t, ok, "Alice account is not found")
	require.NotEmpty(t, alice.Address, "Alice's address is empty")
	require.NotEmpty(t, alice.PrivateKey, "Alice's private key is empty")
	t.Logf("🎉 Extracted Alice: Address: %s, PrivateKey: %s", alice.Address, alice.PrivateKey)

}

func TestSignedTxFromAliceToBob(t *testing.T) {
	client, alice, bob, _ := MustGet(t, devutil.GetUrls())
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
	urls := devutil.GetUrls()
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
	urls := devutil.GetUrls()
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
	urls := devutil.GetUrls()
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

func TestDeployCounterContractViaAPI(t *testing.T) {
	urls := devutil.GetUrls()
	client, alice, _, _ := MustGet(t, urls)
	defer client.Close()

	ctx := context.Background()

	// 🧱 Load contract bytecode from Go binding
	bytecode := counter.CounterMetaData.Bin
	contractAddr, txHash, err := contract.DeployContract(
		ctx, client, urls.ServerURL, alice.Name, bytecode, "",
	)
	require.NoError(t, err)

	t.Logf("🧾 TxHash: %s", txHash)
	t.Logf("🏠 Contract deployed at: %s", contractAddr.Hex())
}

func TestDeployContract_InvalidBytecode(t *testing.T) {
	urls := devutil.GetUrls()
	client, alice, _, _ := MustGet(t, urls)
	defer client.Close()
	ctx := context.Background()

	// 🧱 Get valid bytecode from binding
	bytecode := counter.CounterMetaData.Bin

	// 🧪 Tamper with bytecode: flip some characters near the start
	badBytecode := "0xDEAD" + bytecode[6:]

	contractAddr, _, err := contract.DeployContract(ctx, client, urls.ServerURL, alice.Name, badBytecode, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "contract code is empty")
	t.Logf("✅ Expected error from bad bytecode: %v", err)

	// No need to assert txHash is empty
	require.Equal(t, common.Address{}, contractAddr)
}

func TestRetrieveBlockContentByNumber(t *testing.T) {
	urls := devutil.GetUrls()
	client, _, _, _ := MustGet(t, urls)
	defer client.Close()

	blockNumber := big.NewInt(8)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		logutil.Errorf("❌ Failed to get block by number %d: %v", blockNumber.Int64(), err)
		t.FailNow()
	}

	logutil.Infof("📦 Block #%d", block.Number().Uint64())
	logutil.Infof("🔗 Hash       : %s", block.Hash().Hex())
	logutil.Infof("🔗 Parent Hash: %s", block.ParentHash().Hex())
	logutil.Infof("⛽️ Gas Used   : %d / %d", block.GasUsed(), block.GasLimit())
	logutil.Infof("💥 Transactions: %d", len(block.Transactions()))

	for i, tx := range block.Transactions() {
		logutil.Infof("  ➤ Tx #%d: %s", i, tx.Hash().Hex())

		if tx.To() == nil {
			logutil.Infof("     📦 Contract creation")
		} else {
			logutil.Infof("     📬 To: %s", tx.To().Hex())
		}

		logutil.Infof("     🔢 Nonce: %d | ⛽ Gas: %d | 💰 Value: %s", tx.Nonce(), tx.Gas(), tx.Value().String())
	}
}

func TestDebugTraceTransaction(t *testing.T) {
	urls := devutil.GetUrls()
	client, _, _, _ := MustGet(t, urls)
	defer client.Close()

	txHash := common.HexToHash("0xfe196a1de723b21a066c5d0062d61114059b726b90c835318dfd141bbb9713ed")

	// ⚙️ Manual raw RPC call since `debug_traceTransaction` is not part of ethclient
	var result map[string]interface{}
	rpcClient := client.Client() // This gives us *rpc.Client

	logutil.Infof("🔍 Tracing transaction: %s", txHash.Hex())

	err := rpcClient.CallContext(
		context.Background(),
		&result,
		"debug_traceTransaction",
		txHash,
		map[string]interface{}{}, // default config
	)
	if err != nil {
		logutil.Errorf("❌ Failed to trace transaction: %v", err)
		t.FailNow()
	}

	// 🧠 Print high-level info
	if output, ok := result["output"]; ok {
		logutil.Infof("🧾 Output: %v", output)
	}
	if failed, ok := result["failed"]; ok {
		logutil.Infof("💥 Failed: %v", failed)
	}

	// 🧬 Optional: full dump
	traceBytes, _ := json.MarshalIndent(result, "", "  ")
	logutil.Infof("🧬 Full Trace:\n%s", string(traceBytes))
}

func TestTraceFailedDeployment(t *testing.T) {
	urls := devutil.GetUrls()
	client, _, _, _ := MustGet(t, urls)
	defer client.Close()

	txHash := common.HexToHash("0x25e96caef052cbfb3c24ddf9a5f7f5a2581e7a58082a2a82857db8f51957a7e9")

	var result map[string]interface{}
	rpcClient := client.Client()

	logutil.Infof("🔍 Tracing transaction: %s", txHash.Hex())

	err := rpcClient.CallContext(
		context.Background(),
		&result,
		"debug_traceTransaction",
		txHash,
		map[string]interface{}{}, // default config
	)
	if err != nil {
		logutil.Errorf("❌ Trace failed: %v", err)
		t.FailNow()
	}

	// 🧾 Inspect logs emitted
	if logs, ok := result["structLogs"]; ok {
		logutil.Infof("📜 structLogs present with %d entries", len(logs.([]interface{})))
	}

	traceBytes, _ := json.MarshalIndent(result, "", "  ")
	logutil.Infof("🧬 Full Trace:\n%s", string(traceBytes))
}

func TestGetDeploymentLogs(t *testing.T) {
	urls := devutil.GetUrls()
	client, _, _, _ := MustGet(t, urls)
	defer client.Close()

	// Replace with your deployment tx hash
	txHash := common.HexToHash("0xd9f86c52a716ac4182f8b3629d1f85f4c57950485e6f8e1598202add7a890fa6")

	// First: Get receipt to know block number
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		t.Fatalf("❌ Failed to get receipt: %v", err)
	}
	blockNum := receipt.BlockNumber

	// Second: Filter logs in that block
	query := ethereum.FilterQuery{
		FromBlock: blockNum,
		ToBlock:   blockNum,
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		t.Fatalf("❌ Failed to get logs: %v", err)
	}

	// Third: Display logs
	if len(logs) == 0 {
		logutil.Warnf("🫥 No logs found in block #%d", blockNum.Uint64())
	} else {
		logutil.Infof("📝 Found %d logs in block #%d", len(logs), blockNum.Uint64())
		for i, logEntry := range logs {
			logutil.Infof("🔹 Log %d: Contract=%s", i, logEntry.Address.Hex())
			logutil.Infof("   ➤ Topics: %v", logEntry.Topics)
			logutil.Infof("   ➤ Data  : %x", logEntry.Data)
		}
	}
}
