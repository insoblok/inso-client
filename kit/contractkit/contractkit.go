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
	ABIPath     string // Path to .abi file
	BinPath     string // Path to .bin file
	PackageName string // Name of Go package
	OutputPath  string // Where to write the .go file
}

func GenerateGoBindings(opts BindOptions) error {
	cmd := exec.Command(
		"go", "run", "./cmd/abigen", // 🧠 we’re assuming local go-ethereum clone
		"--abi", opts.ABIPath,
		"--bin", opts.BinPath,
		"--pkg", opts.PackageName,
		"--out", opts.OutputPath,
	)

	cmd.Dir = "/Users/iyadi/github/ethereum/go-ethereum" // 👈 make this configurable if needed

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("abigen failed: %w\nOutput: %s", err, string(out))
	}

	fmt.Printf("✅ abigen success → %s\n", opts.OutputPath)
	return nil
}
