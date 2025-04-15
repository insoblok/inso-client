package types

type SignTxAPIResponse struct {
	SignedTx string `json:"signedTx"`
	TxHash   string `json:"txHash"`
}

type SignTxRequest struct {
	From    string  `json:"from"`              // alias, e.g. "alice"
	To      string  `json:"to"`                // "" or omitted for contract deployment
	Value   string  `json:"value"`             // amount in wei (as string)
	Data    string  `json:"data,omitempty"`    // hex-encoded bytecode or calldata
	Nonce   *uint64 `json:"nonce,omitempty"`   // optional
	ChainID *int64  `json:"chainId,omitempty"` // optional
	Type    string  `json:"type,omitempty"`    // e.g. "deploy", "call", "raw"
}

type DeployContractRequest struct {
	From string `json:"from"`           // alias, e.g. "alice"
	Data string `json:"data,omitempty"` // hex-encoded bytecode or calldata
}

type SendTxAPIResponse struct {
	TxHash string `json:"txHash"`
}
type ContractDeploymentResponse struct {
	TxHash                  string `json:"txHash"`
	ExpectedContractAddress string `json:"expectedContractAddress"`
}

type AliasRegisterResponse struct {
	Status string `json:"status"`
	Alias  string `json:"alias"`
}
