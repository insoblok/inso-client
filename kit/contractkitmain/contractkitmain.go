package main

import (
	"eth-toy-client/kit/contractkit"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: contractkitmain [compile|bind|deploy]")
		os.Exit(1)
	}

	mode := contractkit.Mode(os.Args[1])

	compileOptions := contractkit.CompileOptions{
		SolContractPath: "/Users/iyadi/playground/eth-toy-client/eth-toy-client/sol/contracts/Counter.sol",
		OutBaseDir:      "/Users/iyadi/playground/eth-toy-client/eth-toy-client/sol/out",
		Clean:           true,
	}

	bindOptions := contractkit.BindOptions{
		PackageName: "counter",
		OutFile:     "sol/out/Counter/counter.go",
	}

	//deployOptions := contractkit.DeployOptions{
	//	SolContractPath: "",
	//	OutBaseDir:      "",
	//}

	switch mode {
	case contractkit.ModeCompile:
		fmt.Println("🛠️ Running in COMPILE mode")
		_, err := contractkit.CompileContract(compileOptions)
		if err != nil {
			panic(err)
		}
	case contractkit.ModeBind:
		fmt.Println("🔧 Running in BIND mode")
		_, err := contractkit.RunBind(compileOptions, bindOptions)
		if err != nil {
			panic(err)
		}
	case contractkit.ModeDeploy:
		fmt.Println("🚀 Running in DEPLOY mode")
	default:
		fmt.Printf("❌ Unknown mode: %s\n", mode)
		os.Exit(1)
	}
}
