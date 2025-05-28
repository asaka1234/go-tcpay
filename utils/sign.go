package utils

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"strings"
)

type TCPayRSASignatureUtil struct{}

// category = 1 是: 下单接口
// category = 2 是: 验证接口
func (util *TCPayRSASignatureUtil) GetSign(paramsMap map[string]interface{}, privateKey string, category int) (string, error) {
	delete(paramsMap, "SignData")

	var textContent string
	//1. 拿到签名字符串
	if category == 1 {
		//下单接口
		textContent = util.GetCreatePaymentContent(paramsMap)
	} else if category == 2 {
		//验证接口
		textContent = util.GetVerifyPaymentContent(paramsMap)
	}

	//2. 做sha256 hash
	sha256Data := util.ToSHA256(textContent)
	//3. 用私钥签名
	return util.SignData(2048, privateKey, sha256Data)

}

/*
func (util *TCPayRSASignatureUtil) VerifySign(paramsMap map[string]interface{}, publicKey string, sign string) (bool, error) {
	delete(paramsMap, "platSign")
	textContent := util.GetVerifyContent(paramsMap)
	return util.Verify(textContent, sign, publicKey)
}
*/

//-------------------------------------------------------------

func (util *TCPayRSASignatureUtil) GetCreatePaymentContent(paramsMap map[string]interface{}) string {

	//MerchantId#TerminalId#Action#Amount#InvoiceNumber#LocalDateTime#ReturnUrl#AdditionalData#ConsumerId

	params := ConvertToStringMap(paramsMap)

	var builder strings.Builder
	builder.WriteString(params["MerchantId"])
	builder.WriteString("#")
	builder.WriteString(params["TerminalId"])
	builder.WriteString("#")
	builder.WriteString(params["Action"])
	builder.WriteString("#")
	fmt.Fprintf(&builder, "%.2f#", params["Amount"])
	builder.WriteString(params["InvoiceNumber"])
	builder.WriteString("#")
	builder.WriteString(params["LocalDateTime"])
	builder.WriteString("#")
	builder.WriteString(params["ReturnUrl"])
	builder.WriteString("#")
	builder.WriteString(params["AdditionalData"])
	builder.WriteString("#")
	builder.WriteString(params["ConsumerId"])
	queryString := builder.String()

	return queryString
}

func (util *TCPayRSASignatureUtil) GetVerifyPaymentContent(paramsMap map[string]interface{}) string {
	return cast.ToString(paramsMap["Token"])
}

func (util *TCPayRSASignatureUtil) ToSHA256(input string) string {
	// 创建SHA256哈希对象
	hash := sha256.New()

	// 写入输入数据
	hash.Write([]byte(input))

	// 计算哈希值
	hashBytes := hash.Sum(nil)

	// 将字节数组转换为十六进制字符串
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}

// SignData 使用RSA私钥对数据进行签名
// keySize - 密钥大小（如2048）
// privateKey - PEM格式的私钥
// stringToBeSigned - 要签名的字符串
func (util *TCPayRSASignatureUtil) SignData(keySize int, privateKey string, stringToBeSigned string) (string, error) {
	// 解析PEM格式的私钥
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the private key")
	}

	// 解析私钥
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// 检查密钥大小是否匹配
	if privKey.Size()*8 != keySize {
		return "", fmt.Errorf("key size mismatch: expected %d, got %d", keySize, privKey.Size()*8)
	}

	// 计算SHA256哈希
	hashed := sha256.Sum256([]byte(stringToBeSigned))

	// 使用PKCS#1 v1.5填充进行签名
	signature, err := rsa.SignPKCS1v15(nil, privKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("error signing data: %v", err)
	}

	// 返回Base64编码的签名
	return base64.StdEncoding.EncodeToString(signature), nil
}
