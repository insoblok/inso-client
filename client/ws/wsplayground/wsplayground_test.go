package wsplayground

import "testing"

// Example of using GetAbiBin in tests
func TestGetAbiBin(t *testing.T) {
	logger := &TestLogger{t: t}
	dir := "./path/to/compiled/contracts"
	contract := "MyContract"

	bin, abi := GetAbiBin(logger, dir, contract)
	t.Logf("Binary: %x", bin)
	t.Logf("ABI: %v", abi)
}
