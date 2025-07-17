package pawapaygo

import (
	_ "crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
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
	requestBody, err := payload.ToBytes()
	if err != nil {
		return nil, err
	}

	// Create an http request
	req, err := http.NewRequest("POST", a.instanceURL+"/v2"+requestDepositRoute, requestBody)
	if err != nil {
		return nil, err
	}

	// Add required http headers
	req.Header.Set("Authorization", "Bearer "+a.authToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Signature-Date", time.Now().Format(time.RFC3339))
	digest,err := createContentDigest(payload)
	if err !=nil{
		panic(err)
	}
	params := &SignatureBaseParams{
		DirivedComponents: SignatureBaseDirivedComponents{
			Method: req.Method,
			Authority: req.Host,
			Path: req.URL.Path,
		},
		Headers: SignatureBaseHeaders{
			SignatureDate: req.Header.Get("Signature-Date") ,
			ContentDigest: digest,
			ContentType: req.Header.Get("Content-Type"),
		},
		Metadata: SignatureBaseMetadata{
			Alg: "sha-512",
			KeyId: "somekeyid",
			Created: time.Now().UnixMilli(),
			Expires: time.Now().UnixMilli(),
		},
	}
	fmt.Println(createSignatureBase(params))
	// Make an http request
	res, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Get response body from bytes
	body := &RequestDepositResponse{}
	body.DecodeBytes(res.Body)

	return body, nil
}

func createContentDigest(data any) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	hash := sha512.Sum512(jsonBytes)

	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

func createPrivateKey() {

}

func getPublicKey() string {
	return "private_key"
}

func CreateSignature() string {
	// create signature base
	// use primary to sign the signaturebase
	return "signature"
}

func VerifyPawapaySignature() bool {
	return false
}

func VerifyContentDigest() bool {
	return false
}

func createSignatureBase(params *SignatureBaseParams) string {
	base := []string{
		fmt.Sprintf(`"@method": %s`, params.DirivedComponents.Method),
		fmt.Sprintf(`"@authority": %s`, params.DirivedComponents.Authority),
		fmt.Sprintf(`"@path": %s`, params.DirivedComponents.Path),
		fmt.Sprintf(`"@signature-date": %s`, params.Headers.SignatureDate),
		fmt.Sprintf(`"@content-digest": sha-512=:%s:`, params.Headers.ContentDigest),
		fmt.Sprintf(`"@content-type": %s`, params.Headers.ContentType),
		fmt.Sprintf(`"@signature-params": ("@method" "@authority" "@path" "@signature-date" "@content-digest" "@content-type);alg="%s";keyid="%s";created=%d;expires=%d`, params.Metadata.Alg, params.Metadata.KeyId, params.Metadata.Created, params.Metadata.Expires),
	}

	return strings.Join(base, "\n")
}
