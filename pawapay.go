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
	
	// Create an http request
	req, err := http.NewRequest("POST", a.instanceURL + "/v2" +requestDepositRoute, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Close=true

	// Add required http headers
	req.Header.Set("Authorization", "Bearer "+ a.authToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	// Close request body stream in the end
	defer res.Body.Close()
	
	resBytes,err := io.ReadAll(res.Body)
	if err!=nil{
		return nil, fmt.Errorf("error reading response bytes: %v", err)
	}

	body := &RequestDepositResponse{}
	if err := json.Unmarshal(resBytes, body); err!=nil{
		return nil, err
	}
	
	return body, nil
}
