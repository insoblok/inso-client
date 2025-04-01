package devutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type DevContext struct {
	Client    *ethclient.Client
	ServerURL string
	FromAlias string
	CleanupFn func() // optional if you want to close things later
}

type ClientTestAccount struct {
	Name          string         `json:"name"`
	Address       string         `json:"address"`
	PrivateKey    string         `json:"privateKey"`
	CommonAddress common.Address `json:"-"`
}

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

func GetDevContext(fromAlias string) (*DevContext, error) {
	urls := GetUrls() // local dev server info

	info, err := GetInfoResponse(urls)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch info from devserver: %w", err)
	}

	client, err := ethclient.Dial(info.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ethclient at %s: %w", info.RPCURL, err)
	}

	return &DevContext{
		Client:    client,
		ServerURL: urls.ServerURL,
		FromAlias: fromAlias,
		CleanupFn: func() {
			client.Close()
		},
	}, nil
}
