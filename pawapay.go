package pawapaygo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	instanceURL string
	authToken   string
}

var _ PawapayAPIClient = (*Client)(nil)

func NewPawapayClient(cfg *ConfigOptions) *Client {
	return &Client{
		instanceURL: cfg.InstanceURL,
		authToken:   cfg.ApiToken,
	}
}

type PawapayAPIClient interface {
	InitiateDeposit(*InitiateDepositRequestBody) (*RequestDepositResponse, error)
}

func (a *Client) InitiateDeposit(payload *InitiateDepositRequestBody) (*RequestDepositResponse, error) {

	// Initialize an http client
	httpc := http.Client{}

	// build a request body reader
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(bodyBytes)
	bodyReadCloser := io.NopCloser(bodyReader)
	
	// Create an http request
	req, err := http.NewRequest("POST", a.instanceURL+requestDepositRoute, bodyReadCloser)
	if err != nil {
		return nil, err
	}

	// Add required http headers
	req.Header.Add("Authorization", "Bearer "+ a.authToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	// Close request body stream in the end
	defer res.Body.Close()
	
	var body RequestDepositResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &body, nil
}
