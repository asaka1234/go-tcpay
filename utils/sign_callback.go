package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cast"
	"sort"
	"strings"
)

// 内部自实现的一个md5签名, 用来解决callback没有验签的安全问题

func SignCallback(params map[string]interface{}, key string) (string, error) {

	// 2. Get and sort keys
	keys := []string{"MerchantId", "TerminalId", "Action", "Amount", "InvoiceNumber"}
	sort.Strings(keys) // ASCII ascending order

	// 3. Build sign string
	var sb strings.Builder
	for _, k := range keys {
		value := cast.ToString(params[k])
		if value != "" {
			//只有非空才可以参与签名
			sb.WriteString(fmt.Sprintf("%s=%s&", k, value))
		}
	}
	sb.WriteString(fmt.Sprintf("key=%s", key))
	signStr := sb.String()

	// 4. Generate MD5
	hash := md5.Sum([]byte(signStr))
	signResult := hex.EncodeToString(hash[:])

	// Debug print (optional)
	//fmt.Printf("验签str: %s\n结果: %s\n", signStr, signResult)

	return signResult, nil
}

func VerifyCallback(params map[string]interface{}, signKey string) (bool, error) {
	// Check if signature exists in params
	signature, exists := params["AdditionalData"]
	if !exists {
		return false, nil
	}

	// Remove signature from params for verification
	delete(params, "AdditionalData")

	// Generate current signature
	currentSignature, err := SignCallback(params, signKey)
	if err != nil {
		return false, fmt.Errorf("signature generation failed: %w", err)
	}

	// Compare signatures
	return signature.(string) == currentSignature, nil
}
