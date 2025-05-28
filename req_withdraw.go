package go_tcpay

import (
	"crypto/tls"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

// https://pay-apidoc-en.cheezeebit.com/#p2p-payout-order
func (cli *Client) Withdraw(req CheezeePayWithdrawReq) (*CheezeePayWithdrawResp, error) {

	rawURL := cli.WithdrawURL

	// 1. 拿到请求参数，转为map
	var signDataMap map[string]interface{}
	mapstructure.Decode(req, &signDataMap)
	signDataMap["merchantsId"] = cli.MerchantID
	signDataMap["pushAddress"] = cli.WithdrawCallbackURL
	signDataMap["takerType"] = "2"
	signDataMap["coin"] = "USDT"
	signDataMap["tradeType"] = "1"
	signDataMap["language"] = "en" //TODO 先写死

	// 2. 计算签名,补充参数
	signStr, _ := cli.rsaUtil.GetSign(signDataMap, cli.RSAPrivateKey) //私钥加密
	signDataMap["platSign"] = signStr

	var result CheezeePayWithdrawResp

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

	//验证签名
	if result.Code == "000000" {
		sign := result.PlatSign //收到的签名

		var signResultMap map[string]interface{}
		mapstructure.Decode(result, &signResultMap)
		delete(signResultMap, "platSign") //去掉，用余下的来计算签名

		verify, _ := cli.rsaUtil.VerifySign(signResultMap, cli.RSAPublicKey, sign) //公钥解密
		if !verify {
			return nil, fmt.Errorf("sign verify failed")
		}
	}

	return &result, nil
}
