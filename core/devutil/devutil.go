package devutil

import (
	"context"
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

	// Fetch RPC URL (used to connect ethclient)
	info, err := GetInfoResponse(urls)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch info from devserver: %w", err)
	}

	// Connect to ethclient (RPC endpoint)
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

func GetInfo(ctx context.Context, infoURL string) (*InfoResponse, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, infoURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info InfoResponse
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func GetAccounts(ctx context.Context, url string) ([]ClientTestAccount, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var accounts []ClientTestAccount
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

func makeAliasMap(accounts []ClientTestAccount) map[string]ClientTestAccount {
	m := make(map[string]ClientTestAccount)
	for _, acc := range accounts {
		m[acc.Name] = acc
	}
	return m
}
