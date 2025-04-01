package logutil

import (
	"fmt"
	"os"
)

// ℹ️ Informational log
func Infof(format string, args ...any) {
	fmt.Printf("ℹ️  "+format+"\n", args...)
}

// ⚠️ Warning log
func Warnf(format string, args ...any) {
	fmt.Printf("⚠️  "+format+"\n", args...)
}

// ❌ Error log
func Errorf(format string, args ...any) {
	fmt.Printf("❌ "+format+"\n", args...)
}

// 💀 Fatal error with exit
func Exitf(format string, args ...any) {
	msg := fmt.Sprintf("💀 "+format, args...)
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

// ❌ Create error for propagation (not printed immediately)
func ErrorErrf(format string, args ...any) error {
	return fmt.Errorf("❌ "+format, args...)
}
