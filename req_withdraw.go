package go_tcpay

import (
	"crypto/tls"
	"github.com/asaka1234/go-tcpay/utils"
	"github.com/mitchellh/mapstructure"
	"time"
)

func (cli *Client) Withdraw(req TCPayCreatePaymentReq) (*TCPayCreatePaymentResponse, error) {

	rawURL := cli.CreatePaymentURL

	// 1. 拿到请求参数，转为map
	var signDataMap map[string]interface{}
	mapstructure.Decode(req, &signDataMap)
	signDataMap["MerchantId"] = cli.MerchantID
	signDataMap["TerminalId"] = cli.TerminalID
	signDataMap["LocalDateTime"] = time.Now().Format("2006/01/02 15:04:05")
	signDataMap["Action"] = 100 //100-withdraw
	signDataMap["ReturnUrl"] = cli.WithdrawCallbackURL

	// 2. 先计算一个md5签名, 随后补充到AdditionalData字段中.
	signDataMap["AdditionalData"], _ = utils.SignCallback(signDataMap, cli.RSAPrivateKey)
	
	// 2. 计算签名,补充参数
	signStr, _ := cli.signUtil.GetSign(signDataMap, cli.RSAPrivateKey, 1) //私钥加密
	signDataMap["SignData"] = signStr

	var result TCPayCreatePaymentResponse

	_, err := cli.ryClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetCloseConnection(true).
		R().
		SetBody(signDataMap).
		SetHeaders(getHeaders()).
		SetResult(&result).
		SetError(&result).
		Post(rawURL)

	//fmt.Printf("result: %s\n", string(resp.Body()))

	if err != nil {
		return nil, err
	}

	return &result, nil
}
