package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	pawapay "github.com/Iknite-Space/pawapay-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := &pawapay.ConfigOptions{
		InstanceURL: os.Getenv("BASE_URL"),
		ApiToken:    os.Getenv("AUTH_TOKEN"),
	}
	fmt.Println("CONFIG Variables\n", cfg)
	client := pawapay.NewPawapayClient(cfg)

	router := gin.Default()

	router.POST("/initiate-deposit", func(c *gin.Context) {
		reqBody := &pawapay.InitiateDepositRequestBody{
			DepositID:            uuid.New().String(),
			Amount:               "100",
			Currency:             pawapay.CURRENCY_CODE_CAMEROON,
			PreAuthorisationCode: "54366",
			ClientReferenceID:    "REF-45343",
			CustomerMessage:      "Testing the api",
			Payer: pawapay.Payer{
				Type: "MMO",
				AccountDetails: pawapay.AccountDetails{
					PhoneNumber: "237653456789",
					Provider:    pawapay.MTN_MOMO_CMR,
				},
			},
		}
		res, err := client.InitiateDeposit(reqBody)
		fmt.Println(err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, res)
	})
	router.POST("/deposit-callback", func(c *gin.Context) {
		// fmt.Println(c.Request)
		body := &pawapay.DepositCallbackRequestBody{}
		if err := c.ShouldBindBodyWith(body, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		validated := pawapay.ValidateSignature(c.Request, "keyId", os.Getenv("PrivateKey"))
		fmt.Println(validated)
		c.JSON(http.StatusOK, gin.H{
			"validated": validated,
		})


		for kv := range strings.Split(c.Request.Header.Get("Signature-Input"), ";"){
			if strings.Split(v, "=")[1:]{

			}
		}
		

		sigParams := pawapay.SignatureParams{
			Components: []pawapay.Component{
				{Name: "@method"},
				{Name: "@authority"},
				{Name: "@path"},
				{Name: "content-digest"},
				{Name: "content-type"},
			},
			Created: time.Now().Unix(),
			KeyID:   os.Getenv("KEY_ID"),
			Alg: "sha-512",
		}
		bodyBytes, err := json.Marshal(body)
		if err !=nil{
			fmt.Println("Fail to marshal struct")
		}
		signatureBase, _, err := pawapay.CreateSignatureBase(c.Request, bodyBytes, sigParams)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": err,
			})
		}
		fmt.Println(signatureBase)
		fmt.Println("")
	})
	router.Run()

}