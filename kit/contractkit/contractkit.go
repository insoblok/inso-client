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

func CompileContract(opts CompileOptions) error {
	// 🧠 Derive contract name (no extension)
	base := filepath.Base(opts.SolContractPath)
	name := strings.TrimSuffix(base, filepath.Ext(base)) // "Counter"

	outDir := filepath.Join(opts.OutBaseDir, name) // e.g. "build/Counter"

	// 🧹 Clean if requested
	if opts.Clean {
		if err := os.RemoveAll(outDir); err != nil {
			return fmt.Errorf("failed to clean outdir: %w", err)
		}
	}

	// 📂 Ensure output dir exists
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("failed to create outdir: %w", err)
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
		return fmt.Errorf("solc failed: %w\nOutput: %s", err, string(out))
	}

	fmt.Printf("✅ Compiled %s → %s\n", opts.SolContractPath, outDir)
	return nil
}

type BindOptions struct {
	SolContractPath string
	OutBaseDir      string
}

func BindContract(opts BindOptions) error {
	return nil
}

type DeployOptions struct {
	SolContractPath string
	OutBaseDir      string
}

func DeployContract(opts DeployOptions) error {
	return nil
}
