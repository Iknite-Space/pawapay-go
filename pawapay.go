package pawapaygo

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	hs "github.com/thinkgos/http-signature-go"
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
		fmt.Println("Error converting request body to bytes\n", err)
		return nil, err
	}

	// // Create an http request
	req, err := http.NewRequest("POST", a.instanceURL+"/v2"+requestDepositRoute, requestBody)
	if err != nil {
		fmt.Println("Error creating new request body\n", err)
		return nil, err
	}

	// Add required http headers
	req.Header.Set("Authorization", "Bearer "+a.authToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	res, err := httpc.Do(req)
	if err != nil {
		fmt.Println("Error making an http request to pawapay\n", err)
		return nil, err
	}
	// Close request body stream in the end
	defer res.Body.Close()

	body := &RequestDepositResponse{}
	if err := body.DecodeBytes(res.Body); err != nil {
		fmt.Println("Error parsing the response body to go struct\n", err)
		return nil, err
	}

	return body, nil
}

func ValidateSignature(r *http.Request, keyId string, privateKey string) bool {

	parser := hs.NewParser(
		hs.WithMinimumRequiredHeaders([]string{
			hs.Date,
			hs.Digest,
			hs.HeaderSignature,
			hs.Host,
		}),
		hs.WithSigningMethods(
			hs.SigningMethodRsaPssSha512.Alg(),
			func() hs.SigningMethod { return hs.SigningMethodRsaPssSha512 },
		),
		hs.WithSigningMethods(
			hs.SigningMethodRsaPssSha256.Alg(),
			func() hs.SigningMethod { return hs.SigningMethodRsaPssSha256 },
		),
		// hs.WithValidators(
		// 	hs.NewDigestUsingSharedValidator(),
		// 	hs.NewDateValidator(),
		// ),
		// hs.WithKeystone(keyStone),
	)
	err := parser.AddMetadata(
		hs.KeyId(keyId),
		hs.Metadata{
			Alg:    hs.SigningMethodRsaPssSha512.Name,
			Key:    []byte(""),
			Scheme: hs.SchemeSignature,
		})
	if err != nil {
		fmt.Println(err)
		return false
	}

	gotParam, err := parser.ParseFromRequest(r)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if err := parser.Verify(r, gotParam); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

type Component struct {
	Name       string
	Parameters map[string]string
}

type SignatureParams struct {
	Components []Component
	Alg string
	Created    int64
	KeyID      string
}

// CreateContentDigestHeader generates the SHA-512 content-digest
func CreateContentDigestHeader(body []byte) string {
	sum := sha512.Sum512(body)
	digest := base64.StdEncoding.EncodeToString(sum[:])
	return fmt.Sprintf("sha-512=:%s:", digest)
}

// CreateSignatureBase returns both the signature base string and Signature-Input header value
func CreateSignatureBase(req *http.Request, body []byte, sigParams SignatureParams) (signatureBase, signatureInput string, err error) {
	seen := make(map[string]bool)
	var sb strings.Builder
	var inputNames []string

	// Calculate Content-Digest header and add it if not present
	if req.Header.Get("Content-Digest") == "" {
		digest := CreateContentDigestHeader(body)
		req.Header.Set("Content-Digest", digest)
	}

	for _, comp := range sigParams.Components {
		identifier := serializeComponentIdentifier(comp)
		if seen[identifier] {
			return "", "", fmt.Errorf("duplicate component identifier: %s", identifier)
		}
		seen[identifier] = true
		inputNames = append(inputNames, fmt.Sprintf(`"%s"`, comp.Name))

		// Build signature base line
		sb.WriteString(identifier)
		sb.WriteString(": ")

		value, err := getComponentValue(req, comp)
		if err != nil {
			return "", "", fmt.Errorf("failed to get value for %s: %v", comp.Name, err)
		}
		sb.WriteString(value)
		sb.WriteString("\n")
	}

	// Build final signature-params line
	sigParamsLine := fmt.Sprintf(`"@signature-params": (%s);alg=%s;created=%d;keyid="%s"`,
		strings.Join(inputNames, " "),sigParams.Alg, sigParams.Created, sigParams.KeyID)

	sb.WriteString(sigParamsLine)

	return sb.String(), sigParamsLine, nil
}

func serializeComponentIdentifier(comp Component) string {
	id := fmt.Sprintf(`"%s"`, comp.Name)
	for k, v := range comp.Parameters {
		if v == "" {
			id += ";" + k
		} else {
			id += fmt.Sprintf(`;%s="%s"`, k, v)
		}
	}
	return id
}

func getComponentValue(req *http.Request, comp Component) (string, error) {
	if strings.HasPrefix(comp.Name, "@") {
		switch comp.Name {
		case "@method":
			return req.Method, nil
		case "@authority":
			return req.Host, nil
		case "@path":
			return req.URL.Path, nil
		default:
			return "", fmt.Errorf("unsupported derived component: %s", comp.Name)
		}
	} else {
		values := req.Header[http.CanonicalHeaderKey(comp.Name)]
		if len(values) == 0 {
			return "", fmt.Errorf("header %s not found", comp.Name)
		}
		return strings.Join(values, ", "), nil
	}
}
