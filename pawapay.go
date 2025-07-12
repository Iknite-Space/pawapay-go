package pawapaygo

import (
	"fmt"
)

type Client struct {
	baseUrl string
	apiKey  string
}

var _ PawapayAPIClient = (*Client)(nil)

func NewPawapayClient(cfg *ConfigOptions) *Client {
	return &Client{
		baseUrl: cfg.BaseUrl,
		apiKey:  cfg.ApiKey,
	}
}

type PawapayAPIClient interface {
	RequestDeposit(*RequestDepositBody) (*RequestDepositResponse, error)
}

func (a *Client) RequestDeposit(*RequestDepositBody) (*RequestDepositResponse, error) {
	fmt.Println(a.apiKey, a.baseUrl)
	return nil, nil
}
