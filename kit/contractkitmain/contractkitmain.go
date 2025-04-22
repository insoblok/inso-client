package main

import (
	"eth-toy-client/core/logutil"
	"eth-toy-client/kit/contractkit"
	"fmt"
	"os"
)

func main() {
	logutil.Infof("📣 super duper main started")

	if len(os.Args) < 2 {
		logutil.Exitf("Usage: contractkitmain [compile|bind|deploy]")
	}

	mode := contractkit.Mode(os.Args[1])

	compileOptions := contractkit.CompileOptions{
		SolContractPath: "/Users/iyadi/playground/eth-toy-client/eth-toy-client/sol/contracts/Counter.sol",
		OutBaseDir:      "/Users/iyadi/playground/eth-toy-client/eth-toy-client/sol/out",
		Clean:           true,
	}

	solPath := compileOptions.SolContractPath
	outDir := compileOptions.OutBaseDir

	alias := "alice"

	switch mode {
	case contractkit.ModeCompile:
		fmt.Println("🛠️ Running in COMPILE mode")
		_, err := contractkit.CompileContract(compileOptions)
		if err != nil {
			panic(err)
		}
	case contractkit.ModeBind:
		fmt.Println("🔧 Running in BIND mode")
		_, err := contractkit.RunBind(compileOptions)
		if err != nil {
			panic(err)
		}
	case contractkit.ModeDeploy:
		logutil.Infof("🚀 Running in DEPLOY mode")
		compileOpts := contractkit.CompileOptions{
			SolContractPath: solPath,
			OutBaseDir:      outDir,
			Clean:           true,
		}

		deployOpts := contractkit.DeployOptions{
			FromAlias: alias,
		}
		err := contractkit.RunDeploy(deployOpts, compileOpts)
		if err != nil {
			logutil.Exitf("Deployment failed: %v", err)
		}
	default:
		logutil.Exitf("Usage: contractkitmain [compile|bind|deploy]")
	}
}
