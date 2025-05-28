package utils

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"math/big"
	"strings"
)

type TCPayRSASignatureUtil struct{}

// category = 1 是: 下单接口
// category = 2 是: 验证接口
func (util *TCPayRSASignatureUtil) GetSign(paramsMap map[string]interface{}, privateKeyXML string, category int) (string, error) {
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

	fmt.Printf("=1=>raw: %s\n", textContent)

	//2. 做sha256 hash
	sha256Data := util.ToSHA256(textContent)

	fmt.Printf("=2=>sha256: %s\n", sha256Data)

	//3. 用私钥签名
	signStr, err := util.SignData(2048, privateKeyXML, sha256Data)

	fmt.Printf("=3=>sign: %s\n", signStr)

	return signStr, err

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
	builder.WriteString(decimal.NewFromFloat(cast.ToFloat64(paramsMap["Amount"])).StringFixed(2))
	builder.WriteString("#")
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
func (util *TCPayRSASignatureUtil) SignData(keySize int, privateKeyXML string, stringToBeSigned string) (string, error) {
	// 解析PEM格式的私钥

	privKey := ParseRSAPrivateKeyXML(privateKeyXML)

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

//------------

type RSAKeyValue struct {
	XMLName  xml.Name `xml:"RSAKeyValue"`
	Modulus  string   `xml:"Modulus"`
	Exponent string   `xml:"Exponent"`
	P        string   `xml:"P"`
	Q        string   `xml:"Q"`
	DP       string   `xml:"DP"`
	DQ       string   `xml:"DQ"`
	InverseQ string   `xml:"InverseQ"`
	D        string   `xml:"D"`
}

// 解析xml格式的私钥
func ParseRSAPrivateKeyXML(privateKeyXmlData string) *rsa.PrivateKey {
	var key RSAKeyValue
	err := xml.Unmarshal([]byte(privateKeyXmlData), &key)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse XML: %v", err))
	}

	// 解码所有Base64编码的组件
	modulus, _ := base64.StdEncoding.DecodeString(key.Modulus)
	exponent, _ := base64.StdEncoding.DecodeString(key.Exponent)
	p, _ := base64.StdEncoding.DecodeString(key.P)
	q, _ := base64.StdEncoding.DecodeString(key.Q)
	dp, _ := base64.StdEncoding.DecodeString(key.DP)
	dq, _ := base64.StdEncoding.DecodeString(key.DQ)
	qi, _ := base64.StdEncoding.DecodeString(key.InverseQ)
	d, _ := base64.StdEncoding.DecodeString(key.D)

	// 创建RSA私钥
	privateKey := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: new(big.Int).SetBytes(modulus),
			E: int(new(big.Int).SetBytes(exponent).Int64()),
		},
		D:      new(big.Int).SetBytes(d),
		Primes: []*big.Int{new(big.Int).SetBytes(p), new(big.Int).SetBytes(q)},
		Precomputed: rsa.PrecomputedValues{
			Dp:   new(big.Int).SetBytes(dp),
			Dq:   new(big.Int).SetBytes(dq),
			Qinv: new(big.Int).SetBytes(qi),
		},
	}

	// 验证私钥有效性
	err = privateKey.Validate()
	if err != nil {
		panic(fmt.Sprintf("Invalid private key: %v", err))
	}

	return privateKey
}
