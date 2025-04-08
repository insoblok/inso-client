package wsplayground

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// HexBytes is a custom byte type
type HexBytes []byte

// Logger interface for logging, error handling, and optional exit functionality
type Logger interface {
	Logf(format string, args ...interface{})            // For general logging
	Errorf(format string, args ...interface{})          // For error logging
	FailNow()                                           // To indicate a critical failure and stop (or exit)
	NoError(err error, msg string, args ...interface{}) // Assert no error occurred
}

// TestLogger implementation using *testing.T
type TestLogger struct {
	t *testing.T
}

func (tl *TestLogger) Logf(format string, args ...interface{}) {
	tl.t.Logf(format, args...)
}

func (tl *TestLogger) Errorf(format string, args ...interface{}) {
	tl.t.Errorf(format, args...)
}

func (tl *TestLogger) FailNow() {
	tl.t.FailNow()
}

func (tl *TestLogger) NoError(err error, msg string, args ...interface{}) {
	if err != nil {
		tl.t.Errorf(msg+": %v", append(args, err)...)
		tl.t.FailNow()
	}
}

// ConsoleLogger implementation for console-based logging
type ConsoleLogger struct{}

func (cl *ConsoleLogger) Logf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func (cl *ConsoleLogger) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func (cl *ConsoleLogger) FailNow() {
	fmt.Fprintln(os.Stderr, "❌ Critical error occurred. Exiting...")
	os.Exit(1)
}

func (cl *ConsoleLogger) NoError(err error, msg string, args ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg+": %v\n", append(args, err)...)
		fmt.Fprintln(os.Stderr, "❌ Exiting due to critical error")
		os.Exit(1)
	}
}

// GetAbiBin loads the binary and ABI of a given contract
func GetAbiBin(logger Logger, dir string, contract string) (HexBytes, *abi.ABI) {
	binPath := filepath.Join(dir, contract+".bin")
	rawByte, err := os.ReadFile(binPath)
	logger.NoError(err, "❌ Failed to read compiled .bin file")

	// Ensure proper hex encoding for raw bytes
	if (len(rawByte) % 2) == 1 {
		rawByte = append([]byte("0"), rawByte...)
	}
	destHexByte := make([]byte, len(rawByte)/2)
	_, err = hex.Decode(destHexByte, rawByte)
	logger.NoError(err, "❌ Invalid hex byte code")

	abiPath := filepath.Join(dir, contract+".abi")
	abiData, err := os.ReadFile(abiPath)
	logger.NoError(err, "❌ Failed to read compiled .abi file")

	parsed, err := abi.JSON(strings.NewReader(string(abiData)))
	logger.NoError(err, "❌ Failed to parse .abi file")

	logger.Logf("✅ Successfully loaded ABI and binary for contract: %s", contract)
	return destHexByte, &parsed
}
