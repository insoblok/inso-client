package types

type SignTxAPIResponse struct {
	SignedTx string `json:"signedTx"`
	TxHash   string `json:"txHash"`
}

type SignTxRequest struct {
	From    string  `json:"from"`              // alias, e.g. "alice"
	To      string  `json:"to"`                // recipient address (omit or "" for contract deployment)
	Value   string  `json:"value"`             // amount in wei (as string)
	Data    string  `json:"data,omitempty"`    // hex-encoded contract bytecode or call data
	Nonce   *uint64 `json:"nonce,omitempty"`   // optional
	ChainID *int64  `json:"chainId,omitempty"` // optional
}

type SendTxAPIResponse struct {
	TxHash string `json:"txHash"`
}

type AliasRegisterResponse struct {
	Status string `json:"status"`
	Alias  string `json:"alias"`
}
