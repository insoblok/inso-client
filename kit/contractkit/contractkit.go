package contractkit

import (
	"context"
	contract "eth-toy-client/core/contracts"
	"eth-toy-client/core/devutil"
	"eth-toy-client/core/logutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Mode determines what task to run
type Mode string

const (
	ModeCompile Mode = "COMPILE"
	ModeBind    Mode = "BIND"
	ModeDeploy  Mode = "DEPLOY"
)

// CompileOptions holds inputs for the compiler step
type CompileOptions struct {
	SolContractPath string // Full path to .sol file
	OutBaseDir      string // Absolute base directory
	Clean           bool   // If true, delete output dir before build
}

// BindOptions holds inputs for the abigen step
type BindOptions struct {
	PackageName string // Go package name for generated code
	OutFile     string // Path for generated .go file
}

// DeployOptions holds runtime info for deployment
type DeployOptions struct {
	RPCURL    string // JSON-RPC endpoint
	FromAlias string // Sender account alias (e.g., "alice")
	ServerURL string // DevServer API endpoint (e.g., http://localhost:8575)
}

// BuildResult represents outputs of the compile step
type BuildResult struct {
	BuildDir     string
	ContractName string
}

func CompileContract(opts CompileOptions) (*BuildResult, error) {
	if _, err := os.Stat(opts.SolContractPath); os.IsNotExist(err) {
		return nil, logutil.ErrorErrf("sol contract file does not exist: %s", opts.SolContractPath)
	}

	if _, err := os.Stat(opts.OutBaseDir); os.IsNotExist(err) {
		logutil.Infof("Output base dir not found — creating: %s", opts.OutBaseDir)
		if err := os.MkdirAll(opts.OutBaseDir, 0o755); err != nil {
			return nil, logutil.ErrorErrf("failed to create OutBaseDir: %w", err)
		}
		logutil.Infof("✅ Created OutBaseDir: %s", opts.OutBaseDir)
	}

	contractName := strings.TrimSuffix(
		filepath.Base(opts.SolContractPath),
		filepath.Ext(opts.SolContractPath),
	)

	logutil.Infof(
		"Compiling contract: contract=%s solContractPath=%s outBaseDir=%s clean=%v",
		contractName, opts.SolContractPath, opts.OutBaseDir, opts.Clean)

	buildDir := filepath.Join(opts.OutBaseDir, contractName)

	if opts.Clean {
		_ = os.RemoveAll(buildDir)
	}

	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return nil, logutil.ErrorErrf("failed to create build dir: %w", err)
	}

	cmd := exec.Command(
		"solc",
		"--abi",
		"--bin",
		"-o",
		buildDir,
		opts.SolContractPath,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, logutil.ErrorErrf("solc failed: %w\nOutput: %s", err, string(out))
	}

	logutil.Infof("✅ Compiled %s → %s", opts.SolContractPath, buildDir)

	return &BuildResult{
		BuildDir:     buildDir,
		ContractName: contractName,
	}, nil
}

func RunBind(compileOpts CompileOptions, bindOpts BindOptions) (*BuildResult, error) {
	result, err := CompileContract(compileOpts)
	if err != nil {
		return nil, logutil.ErrorErrf("compilation failed before binding: %w", err)
	}

	abiFile := filepath.Join(result.BuildDir, result.ContractName+".abi")
	binFile := filepath.Join(result.BuildDir, result.ContractName+".bin")
	goFile := filepath.Join(result.BuildDir, strings.ToLower(result.ContractName)+".go")
	packageName := strings.ToLower(result.ContractName)

	logutil.Infof("ABI file: %s", abiFile)
	logutil.Infof("BIN file: %s", binFile)
	logutil.Infof("Go file: %s", goFile)

	cmd := exec.Command(
		"abigen",
		"--abi", abiFile,
		"--bin", binFile,
		"--pkg", packageName,
		"--out", goFile,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, logutil.ErrorErrf("abigen failed: %w\nOutput: %s", err, string(out))
	}

	logutil.Infof("✅ abigen: %s → %s", abiFile, bindOpts.OutFile)
	return result, nil
}

func RunDeploy(opts DeployOptions, compileOpts CompileOptions, bindOpts BindOptions) error {
	logutil.Infof("🚀 Deploying contract from alias: %s", opts.FromAlias)

	devCtx, err := devutil.GetDevContext(opts.FromAlias)
	if err != nil {
		return logutil.ErrorErrf("failed to setup dev context: %w", err)
	}
	defer devCtx.Client.Close()

	result, err := RunBind(compileOpts, bindOpts)
	if err != nil {
		return logutil.ErrorErrf("compilation failed: %w", err)
	}

	binPath := filepath.Join(result.BuildDir, result.ContractName+".bin")
	bytecode, err := os.ReadFile(binPath)
	if err != nil {
		return logutil.ErrorErrf("failed to read bin file: %w", err)
	}

	addr, txHash, err := contract.DeployContract(context.Background(), devCtx.Client, devCtx.ServerURL, devCtx.FromAlias, "0x"+string(bytecode))
	if err != nil {
		return logutil.ErrorErrf("contract deployment failed: %w", err)
	}

	logutil.Infof("✅ Contract deployed at: %s (tx: %s)", addr.Hex(), txHash)
	return nil
}
