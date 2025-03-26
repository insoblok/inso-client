package testkeys

import (
	"github.com/ethereum/go-ethereum/common"
)

type TestAccount struct {
	Name    string
	PrivHex string
	Addr    common.Address
}

var TestAccounts = []TestAccount{
	{
		Name:    "Alpha",
		PrivHex: "4f3edf983ac636a65a842ce7c78d9aa706d3b113b37d80c7e0e4b3e8d3e4e8f4",
		Addr:    common.HexToAddress("0x90f8bf6a479f320ead074411a4b0e7944ea8c9c1"),
	},
	{
		Name:    "Beta",
		PrivHex: "6c8756c099f118313a0955c53d067ee0f8435b7e56e2a4d5e19df2c207da244d",
		Addr:    common.HexToAddress("0xffcf8fdee72ac11b5c542428b35eef5769c409f0"),
	},
	{
		Name:    "Gamma",
		PrivHex: "646f1ce2c7e7b67d172d775f9500bbf509a2287f657b6a3e6b13fdc7c72b5d36",
		Addr:    common.HexToAddress("0x627306090abaB3A6e1400e9345bC60c78a8BEf57"),
	},
	{
		Name:    "Delta",
		PrivHex: "add53f9a7e588d27c4d13c43e64c8bd67160c2e1f6e0298c0a271e3c8e51ef17",
		Addr:    common.HexToAddress("0xf17f52151EbEF6C7334FAD080c5704D77216b732"),
	},
	{
		Name:    "Epsilon",
		PrivHex: "c88b703fb08cbea894b30858dfdd2d48ff7cc3e94b39c1302c6ea4ce00c9f2c3",
		Addr:    common.HexToAddress("0x5AEDA56215b167893e80B4fE645BA6d5Bab767DE"),
	},
}
