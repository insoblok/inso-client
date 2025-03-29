package types

type SignTxAPIResponse struct {
	SignedTx string `json:"signedTx"`
	TxHash   string `json:"txHash"`
}

type SignTxRequest struct {
	From    string  `json:"from"`
	To      string  `json:"to"`
	Value   string  `json:"value"`             // optional (wei, as string)
	Nonce   *uint64 `json:"nonce,omitempty"`   // optional
	ChainID *int64  `json:"chainId,omitempty"` // optional
}

type SendTxAPIResponse struct {
	TxHash string `json:"txHash"`
}
