package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	pawapay "github.com/iknite-space/pawapay-go"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cfg := &pawapay.ConfigOptions{
		InstanceURL: os.Getenv("BASE_URL"),
		ApiToken:    os.Getenv("AUTH_TOKEN"),
	}

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
				PhoneNumber: "237653456019",
				Provider:    pawapay.MTN_MOMO_CMR,
			},
		},
	}

	client := pawapay.NewPawapayClient(cfg)
	res, err := client.InitiateDeposit(reqBody)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(res)
}
