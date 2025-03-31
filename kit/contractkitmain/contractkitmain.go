package main

import "eth-toy-client/kit/contractkit"

func main() {
	options := contractkit.CompileOptions{
		SolContractPath: "/Users/iyadi/playground/eth-toy-client/eth-toy-client/sol/contracts/Counter.sol",
		OutBaseDir:      "/Users/iyadi/playground/eth-toy-client/eth-toy-client/sol/out",
		Clean:           true,
	}

	err := contractkit.CompileContract(options)
	if err != nil {
		panic(err)
	}

}
