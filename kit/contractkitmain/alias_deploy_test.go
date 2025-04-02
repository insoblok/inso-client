package main

import (
	"eth-toy-client/core/logutil"
	"eth-toy-client/kit/contractkit"
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
