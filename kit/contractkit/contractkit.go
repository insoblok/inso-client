package contractkit

import (
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
	BuildDir string // Relative to OutBaseDir
	ABIPath  string // Relative to BuildDir
	BINPath  string // Relative to BuildDir
	Contract string // Contract name (no extension)
}

// CompileContract compiles the contract and returns build metadata
func CompileContract(opts CompileOptions) (*BuildResult, error) {
	// ✅ Validate contract path
	if _, err := os.Stat(opts.SolContractPath); os.IsNotExist(err) {
		return nil, logutil.ErrorErrf("sol contract file does not exist: %s", opts.SolContractPath)
	}

	// 📁 Ensure output base dir exists
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

	// ✅ Create output dir if needed
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return nil, logutil.ErrorErrf("failed to create build dir: %w", err)
	}

	// ✅ Run solc
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
		return nil, logutil.ErrorErrf("solc failed: %w\nOutput: %s", err, string(out))
	}

	logutil.Infof("✅ Compiled %s → %s", opts.SolContractPath, buildDir)

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
		return nil, logutil.ErrorErrf("compilation failed before binding: %w", err)
	}

	abiFile := filepath.Join(compileOpts.OutBaseDir, result.BuildDir, result.ABIPath)
	binFile := filepath.Join(compileOpts.OutBaseDir, result.BuildDir, result.BINPath)

	logutil.Infof("ABI file: %s\n", abiFile)
	logutil.Infof("BIN file: %s\n", binFile)

	cmd := exec.Command(
		"abigen",
		"--abi", abiFile,
		"--bin", binFile,
		"--pkg", bindOpts.PackageName,
		"--out", bindOpts.OutFile,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, logutil.ErrorErrf("abigen failed: %w\nOutput: %s", err, string(out))
	}

	logutil.Infof("✅ abigen: %s → %s\n", abiFile, bindOpts.OutFile)
	return result, nil
}
