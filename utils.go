package pawapaygo

func (c *Client) validateRequestSignature(signature string) bool {

	return signature != ""
}

