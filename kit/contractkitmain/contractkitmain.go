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

	bindOptions := contractkit.BindOptions{
		PackageName: "counter",
		OutFile:     "sol/out/Counter/counter.go",
	}

	//// Compiler Options
	//solPath := flag.String("sol", "", "Path to Solidity contract")
	//outDir := flag.String("out", "", "Base output directory")
	//clean := flag.Bool("clean", false, "Clean build directory before compile")
	//
	//// Bind Options
	//pkg := flag.String("pkg", "", "Go package name for binding")
	//bindOut := flag.String("bindout", "", "Go file to write binding")
	//
	//// Compiler Options
	solPath := compileOptions.SolContractPath
	outDir := compileOptions.OutBaseDir

	// Bind Options
	pkg := bindOptions.PackageName
	bindOut := bindOptions.OutFile
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
		_, err := contractkit.RunBind(compileOptions, bindOptions)
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
		bindOpts := contractkit.BindOptions{
			PackageName: pkg,
			OutFile:     bindOut,
		}
		deployOpts := contractkit.DeployOptions{
			FromAlias: alias,
		}
		err := contractkit.RunDeploy(deployOpts, compileOpts, bindOpts)
		if err != nil {
			logutil.Exitf("Deployment failed: %v", err)
		}
	default:
		logutil.Exitf("Usage: contractkitmain [compile|bind|deploy]")
	}
}
