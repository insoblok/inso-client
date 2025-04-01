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
	Client   *ethclient.Client
	URLs     Urls
	Accounts map[string]ClientTestAccount
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

// GetDevContext connects to the DevServer and returns RPC + accounts
func GetDevContext(ctx context.Context) (*DevContext, error) {
	urls := GetUrls()

	info, err := GetInfo(ctx, urls.InfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch /info: %w", err)
	}

	client, err := ethclient.DialContext(ctx, info.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	accounts, err := GetAccounts(ctx, urls.AccountsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch /accounts: %w", err)
	}

	// Parse addresses
	for i := range accounts {
		accounts[i].CommonAddress = common.HexToAddress(accounts[i].Address)
	}

	return &DevContext{
		Client:   client,
		URLs:     urls,
		Accounts: makeAliasMap(accounts),
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
