package main

import (
	"fmt"

	pawapay "github.com/iknite-space/pawapay-go"
)

func main() {
	 client := pawapay.NewDefaultClient()
	res, err := client.RequestDeposit(&pawapay.RequestDepositBody{})
	if err != nil {
		fmt.Printf("Error occured while creating a payment request: %v", err)
	}
	fmt.Println(res)
}
