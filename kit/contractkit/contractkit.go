package contractkit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CompileOptions struct {
	SolContractPath string // e.g. "contracts/Counter.sol"
	OutBaseDir      string // e.g. "build/"
	Clean           bool   // clean target dir before compiling
}

func CompileContract(opts CompileOptions) (string, error) {
	// 🧠 Derive contract name (no extension)
	base := filepath.Base(opts.SolContractPath)
	name := strings.TrimSuffix(base, filepath.Ext(base)) // "Counter"

	outDir := filepath.Join(opts.OutBaseDir, name) // e.g. "build/Counter"

	// 🧹 Clean if requested
	if opts.Clean {
		if err := os.RemoveAll(outDir); err != nil {
			return "", fmt.Errorf("failed to clean outdir: %w", err)
		}
	}

	// 📂 Ensure output dir exists
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create outdir: %w", err)
	}

	// 🛠️ Run solc
	cmd := exec.Command(
		"solc",
		"--abi",
		"--bin",
		"-o", outDir,
		opts.SolContractPath,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("solc failed: %w\nOutput: %s", err, string(out))
	}

	fmt.Printf("✅ Compiled %s → %s\n", opts.SolContractPath, outDir)
	return outDir, nil
}

type BindOptions struct {
	PackageName string // Go package name for the bindings
	OutFile     string // output Go file path (e.g. bindings/counter.go)
}

func RunBind(compileOpts CompileOptions, bindOpts BindOptions) error {
	outBaseDir, err := CompileContract(compileOpts)
	if err != nil {
		return fmt.Errorf("compilation failed before binding: %w", err)
	}

	contractDir := outBaseDir
	contractName := filepath.Base(contractDir)

	abiFile := filepath.Join(contractDir, contractName+".abi")
	binFile := filepath.Join(contractDir, contractName+".bin")

	fmt.Println("ABI file:", abiFile)
	fmt.Println("BIN file:", binFile)

	// 🧬 Run abigen
	cmd := exec.Command(
		"abigen",
		"--abi", abiFile,
		"--bin", binFile,
		"--pkg", bindOpts.PackageName,
		"--out", bindOpts.OutFile,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("abigen failed: %w\nOutput: %s", err, string(out))
	}

	fmt.Printf("✅ abigen: %s → %s\n", abiFile, bindOpts.OutFile)
	return nil
}

type DeployOptions struct {
	SolContractPath string
	OutBaseDir      string
}

func DeployContract(opts DeployOptions) error {
	return nil
}
