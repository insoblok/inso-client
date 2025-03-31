package contractkit

import (
	"fmt"
	"os/exec"
)

type CompileOptions struct {
	SolContractPath string // Path to Solidity source (e.g. contracts/Counter.sol)
	OutDir          string // Output directory for .abi and .bin
}

// CompileContract runs solc to generate ABI and BIN files for the given Solidity contract
func CompileContract(opts CompileOptions) error {
	cmd := exec.Command(
		"solc",
		"--abi",
		"--bin",
		"--overwrite",
		"-o", opts.OutDir,
		opts.SolContractPath,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("solc failed: %w\nOutput: %s", err, string(out))
	}

	fmt.Printf("✅ Compiled %s → %s\n", opts.SolContractPath, opts.OutDir)
	return nil
}
