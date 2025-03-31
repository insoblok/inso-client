package contractkit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CompileOptions struct {
	SolContractPath string // Path to Solidity source (e.g. contracts/Counter.sol)
	OutDir          string // Output directory for .abi and .bin
}

func CompileContract(opts CompileOptions) error {
	contractName := strings.TrimSuffix(filepath.Base(opts.SolContractPath), ".sol")
	outDir := filepath.Join(opts.OutDir, contractName)

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	cmd := exec.Command(
		"solc",
		"--abi",
		"--bin",
		"--overwrite",
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
