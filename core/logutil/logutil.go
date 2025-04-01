package logutil

import (
	"fmt"
	"os"
)

func Info(msg string, args ...any) {
	fmt.Printf("ℹ️  "+msg+"\n", args...)
}

func Warn(msg string, args ...any) {
	fmt.Printf("⚠️  "+msg+"\n", args...)
}

func Error(msg string, args ...any) {
	fmt.Printf("❌ "+msg+"\n", args...)
}

// LogAndExit prints a message and exits the program
func LogAndExit(msg string) {
	_, err := fmt.Fprintf(os.Stderr, "❌ %s\n", msg)
	if err != nil {
		return
	}
	os.Exit(1)
}

// Errorf returns a formatted error with ❌ prefix
func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf("❌ "+format, args...)
}
