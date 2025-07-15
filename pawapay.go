package pawapaygo

import (
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

	// // build a request body reader
	// bodyBytes, err := json.Marshal(payload)
	// if err != nil {
	// 	return nil, err
	// }
	requestBody, err := payload.ToBytes()
	if err != nil{
		return nil, err
	}
	
	// // Create an http request
	req, err := http.NewRequest("POST", a.instanceURL + "/v2" +requestDepositRoute, requestBody)
	if err != nil {
		return nil, err
	}

	// Add required http headers
	req.Header.Set("Authorization", "Bearer "+ a.authToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	// Close request body stream in the end
	defer res.Body.Close()
	
	body := &RequestDepositResponse{}
	body.DecodeBytes(res.Body)
	
	return body, nil
}
