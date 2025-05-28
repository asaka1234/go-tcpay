package go_tcpay

import (
	"crypto/tls"
	"fmt"
	"github.com/asaka1234/go-tcpay/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

func (cli *Client) DepositCallback(req TCPayCreatePaymentBackReq, processor func(TCPayCreatePaymentBackReq) error) error {
	//1. 验证签名
	if req.ResCode == 0 {
		var signDataMap map[string]interface{}
		mapstructure.Decode(req.Data, &signDataMap)

		verifyResult, err := utils.VerifyCallback(signDataMap, cli.RSAPrivateKey)
		if err != nil || !verifyResult {
			return fmt.Errorf("illegal callback!")
		}

		if cast.ToString(req.Data.MerchantId) != cli.MerchantID || cast.ToString(req.Data.TerminalId) != cli.TerminalID {
			return fmt.Errorf("illegal merchantID!")
		}

		if req.Data.Action != 50 {
			return fmt.Errorf("illegal action!")
		}
	}

	//2. 业务侧开始处理
	err := processor(req)
	if err != nil {
		return err
	}

	//3. 最后去confirm
	verifyReq := TCPayVerifyPaymentReq{
		Token: req.Data.Token,
	}
	return cli.VerifyPayment(verifyReq)
}

// 收到callback后15分钟内必须发送
func (cli *Client) VerifyPayment(req TCPayVerifyPaymentReq) error {

	var signDataMap map[string]interface{}
	mapstructure.Decode(req, &signDataMap)

	signStr, _ := cli.signUtil.GetSign(signDataMap, cli.RSAPrivateKey, 2) //私钥加密
	signDataMap["SignData"] = signStr

	//---------------------------------

	var result TCPayVerifyPaymentResp

	resp, err := cli.ryClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetCloseConnection(true).
		R().
		SetBody(signDataMap).
		SetHeaders(getHeaders()).
		SetResult(&result).
		SetError(&result).
		Post(cli.VerifyPaymentURL)

	fmt.Printf("verify result: %s, %+v\n", string(resp.Body()), result)

	return err
}
