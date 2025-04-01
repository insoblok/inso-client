package devutil

import (
	"encoding/json"
	"io"
	"net/http"
)

// devutil/info.go

type InfoResponse struct {
	RPCURL        string `json:"rpcUrl"`
	AccountsCount int    `json:"accountsCount"`
}

type Urls struct {
	ServerURL   string
	InfoURL     string
	AccountsURL string
}

func GetUrls() Urls {
	base := "http://localhost:8575" // or env override in future
	return Urls{
		ServerURL:   base,
		InfoURL:     base + "/info",
		AccountsURL: base + "/accounts",
	}
}

func GetInfoResponse(urls Urls) (*InfoResponse, error) {
	resp, err := http.Get(urls.InfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info InfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
