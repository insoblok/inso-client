package main

import (
	"eth-toy-client/kit/contractkit"
	"fmt"
	"os"
)

// Mode defines what operation we want to run
type Mode string

const (
	ModeCompile Mode = "compile"
	ModeBind    Mode = "bind"
	ModeDeploy  Mode = "deploy"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: contractkitmain [compile|bind|deploy]")
		os.Exit(1)
	}

	mode := Mode(os.Args[1])

	switch mode {
	case ModeCompile:
		options := contractkit.CompileOptions{
			SolContractPath: "/Users/iyadi/playground/eth-toy-client/eth-toy-client/sol/contracts/Counter.sol",
			OutBaseDir:      "/Users/iyadi/playground/eth-toy-client/eth-toy-client/sol/out",
			Clean:           true,
		}
		err := contractkit.CompileContract(options)
		if err != nil {
			panic(err)
		}

		fmt.Println("🛠️ Running in COMPILE mode")
	case ModeBind:
		fmt.Println("🔧 Running in BIND mode")
	case ModeDeploy:
		fmt.Println("🚀 Running in DEPLOY mode")
	default:
		fmt.Printf("❌ Unknown mode: %s\n", mode)
		os.Exit(1)
	}
}
