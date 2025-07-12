package main

import (
	"fmt"

	pawapay "github.com/iknite-space/pawapay-go"
)

func main() {
	client := pawapay.NewPawapayClient(&pawapay.ConfigOptions{
		// e.g, https://api.sandbox.pawapay.io
		InstanceURL: "https://dashboard.sandbox.pawapay.io",

		// e.g, 3c0cb67a-c85a-2f6b-b944-39de61d67312
		ApiToken: "3c0cb67a-c84a-4f6b-b844-39de61d67311",
	})
	res, err := client.InitiateDeposit(&pawapay.InitiateDepositRequestBody{
		DepositID: "c768f57f-92ed-4fb3-9247-db01ff994cc4",
		Amount:    "500",
		Currency:  "CMR",
		Payer: pawapay.Payer{
			Type: "MMO", AccountDetails: pawapay.AccountDetails{
				PhoneNumber: "237653456019",
				Provider:    "MTN_MOMO_CMR",
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(res.FailureReason.FailureCode)
}
