package main

import (
	"context"
	"eth-toy-client/core/logutil"
	"eth-toy-client/kit/contractkit"
	"eth-toy-client/kit/mockusdc"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"testing"
)

const (
	SolBasePath = "/Users/iyadi/playground/eth-toy-client/eth-toy-client/sol/contracts/"
	OutBasePath = "/Users/iyadi/playground/eth-toy-client/eth-toy-client/sol/out/"
	Counter     = "Counter"
	CounterV1   = "CounterV1"
	MockUSDC    = "MockUSDC"
	MockUSDCV1  = "MockUSDCV1"
)

var ContractsMap = map[string]string{
	CounterV1:  Contract("Counter").ContractPath(),
	MockUSDCV1: Contract("MockUSDC").ContractPath(),
}

type Contract string

func (contractName Contract) ContractPath() string {
	return SolBasePath + string(contractName) + ".sol"
}

func deploySingleContract(contractName string, owner string, alias string) {
	compileOptions := contractkit.CompileOptions{
		SolContractPath: Contract(contractName).ContractPath(),
		OutBaseDir:      OutBasePath,
		Clean:           true,
	}
	deployOpts := contractkit.DeployOptions{
		FromAlias: owner,
	}

	err := contractkit.RunAliasDeploy(alias, compileOptions, deployOpts)
	if err != nil {
		logutil.Exitf("Deployment failed: %v", err)
	}
}

///////////////////////

func TestComposeContractPath(t *testing.T) {
	counter := Contract("Counter")
	logutil.Infof("Counter path: %s", counter.ContractPath())
}

func TestALiasDeploySingleContract(t *testing.T) {
	owner := "alice"
	//deploySingleContract("IntTestExploder", owner, "IntTestExploder")
	//deploySingleContract("CounterV2", owner, "CounterV2")
	deploySingleContract(MockUSDC, owner, MockUSDCV1)
	//deploySingleContract(Counter, owner, CounterV1)
}

func TestALiasDeployContracts(t *testing.T) {
	alias := "alice"

	for contractKey, contractPath := range ContractsMap {
		logutil.Infof("Deploying Contract : %s, Path: %s", contractKey, contractPath)
		compileOptions := contractkit.CompileOptions{
			SolContractPath: contractPath,
			OutBaseDir:      OutBasePath,
			Clean:           true,
		}
		deployOpts := contractkit.DeployOptions{
			FromAlias: alias,
		}
		err := contractkit.RunAliasDeploy(contractKey, compileOptions, deployOpts)
		if err != nil {
			logutil.Warnf("Deployment failed: %v Contract: %s", err, contractKey)
		}
	}

}

//func TestBindDeploy(t *testing.T) {
//	chainID := big.NewInt(1337)
//	//alias := "alice"
//	key := "0x4f3edf983ac636a65a842ce7c78d9aa706d3b113b37b1e1a7e3975d3fef1fdc3"
//	privKey, err := crypto.HexToECDSA(key[2:])
//	if err != nil {
//		log.Fatalf("❌ failed to get private key", err)
//	}
//	address := crypto.PubkeyToAddress(privKey.PublicKey)
//	log.Printf("address: %v", address.Hex())
//	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Set up tx options (you can tweak gas settings if needed)
//	auth.GasLimit = uint64(300000)   // or estimate
//	auth.GasFeeCap = big.NewInt(1e9) // optional
//	auth.GasTipCap = big.NewInt(1)   // optional
//
//	// deploy using generated code
//	backend := ""
//	address, tx, instance, err := mockusdc.DeployMockusdc(auth, backend)
//	if err != nil {
//		log.Fatalf("💥 Failed to deploy: %v", err)
//	}
//
//	fmt.Println("📦 Contract deployed at:", address.Hex())
//	fmt.Println("📄 TX hash:", tx.Hash().Hex())
//}

func TestBindDeployUSDC(t *testing.T) {
	client, err := ethclient.Dial("http://127.0.0.1:8565")
	if err != nil {
		log.Fatalf("❌ Failed to connect to node: %v", err)
	}
	defer client.Close()

	//key := "0x4f3edf983ac636a65a842ce7c78d9aa706d3b113b37b1e1a7e3975d3fef1fdc3"
	//// 2. Load your private key (from funded dev account)
	//privateKeyHex := key[:2] // ⛔ keep secret!
	privateKey, err := crypto.HexToECDSA("4f3edf983ac636a65a842ce7c78d9aa706d3b113b37b1e1a7e3975d3fef1fdc3")
	if err != nil {
		log.Fatalf("❌ Failed to load private key: %v", err)
	}

	// 3. Prepare auth transactor
	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	chainID, _ := client.ChainID(context.Background())
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("❌ Failed to create transactor: %v", err)
	}

	// Optional: set nonce/gas manually
	nonce, _ := client.PendingNonceAt(context.Background(), fromAddr)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // no ETH to send
	auth.GasLimit = uint64(3000000) // deployment cost
	auth.GasPrice = big.NewInt(1e9)

	address, tx, instance, err := mockusdc.DeployMockusdc(auth, client)
	if err != nil {
		log.Fatalf("💥 Deployment failed: %v", err)
	}

	fmt.Println("📦 Contract instance:", instance)
	fmt.Println("📦 Contract deployed at:", address.Hex())
	fmt.Println("📄 TX hash:", tx.Hash().Hex())

}
