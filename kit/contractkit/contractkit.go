package contractkit

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
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
	BuildDir string // Relative to OutBaseDir
	ABIPath  string // Relative to BuildDir
	BINPath  string // Relative to BuildDir
	Contract string // Contract name (no extension)
}

// CompileContract compiles the contract and returns build metadata
func CompileContract(opts CompileOptions) (*BuildResult, error) {
	contractName := strings.TrimSuffix(
		filepath.Base(opts.SolContractPath),
		filepath.Ext(opts.SolContractPath),
	)

	log.Info(
		"Compiling contract",
		"contract", contractName,
		"solContractPath", opts.SolContractPath,
		"outBaseDir", opts.OutBaseDir,
		"clean", opts.Clean)

	buildDir := filepath.Join(opts.OutBaseDir, contractName)

	if opts.Clean {
		_ = os.RemoveAll(buildDir)
	}

	// Create output dir if not exists
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create build dir: %w", err)
	}

	cmd := exec.Command(
		"solc",
		"--abi",
		"--bin",
		"--overwrite",
		"-o", buildDir,
		opts.SolContractPath,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("solc failed: %w\nOutput: %s", err, string(out))
	}

	fmt.Printf("✅ Compiled %s → %s\n", opts.SolContractPath, buildDir)

	return &BuildResult{
		BuildDir: buildDir,
		ABIPath:  filepath.Base(contractName + ".abi"),
		BINPath:  filepath.Base(contractName + ".bin"),
		Contract: contractName,
	}, nil
}

// RunBind compiles and binds the contract
func RunBind(compileOpts CompileOptions, bindOpts BindOptions) (*BuildResult, error) {
	result, err := CompileContract(compileOpts)
	if err != nil {
		return nil, fmt.Errorf("compilation failed before binding: %w", err)
	}

	abiFile := filepath.Join(compileOpts.OutBaseDir, result.BuildDir, result.ABIPath)
	binFile := filepath.Join(compileOpts.OutBaseDir, result.BuildDir, result.BINPath)

	fmt.Printf("ABI file: %s\n", abiFile)
	fmt.Printf("BIN file: %s\n", binFile)

	cmd := exec.Command(
		"abigen",
		"--abi", abiFile,
		"--bin", binFile,
		"--pkg", bindOpts.PackageName,
		"--out", bindOpts.OutFile,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("abigen failed: %w\nOutput: %s", err, string(out))
	}

	fmt.Printf("✅ abigen: %s → %s\n", abiFile, bindOpts.OutFile)
	return result, nil
}
