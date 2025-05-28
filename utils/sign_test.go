package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestSign(t *testing.T) {
	// Example usage
	privateKey := `-----BEGIN RSA PRIVATE KEY-----
... your private key here ...
-----END RSA PRIVATE KEY-----`

	publicKey := `-----BEGIN PUBLIC KEY-----
... your public key here ...
-----END PUBLIC KEY-----`

	params := map[string]interface{}{
		"amount":   "100.00",
		"currency": "USD",
		"orderId":  "123456",
		"data": map[string]interface{}{
			"userId": "user123",
		},
	}

	// Generate signature
	signature, err := GetSign(params, privateKey)
	if err != nil {
		log.Fatal("Sign error:", err)
	}
	fmt.Println("Generated signature:", signature)

	// Verify signature
	params["platSign"] = signature
	valid, err := VerifySign(params, publicKey)
	if err != nil {
		log.Fatal("Verify error:", err)
	}
	fmt.Println("Signature valid:", valid)

}
