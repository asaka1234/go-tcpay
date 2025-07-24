package go_tcpay

import (
	"crypto/tls"
	"fmt"
	"github.com/asaka1234/go-tcpay/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

func (cli *Client) DepositCallback(req TCPayCreatePaymentBackReq, processor func(TCPayCreatePaymentBackReq) error) error {
	//1. 验证签名
	if req.ResCode == 0 {
		var signDataMap map[string]interface{}
		mapstructure.Decode(req.Data, &signDataMap)

		verifyResult, err := utils.VerifyCallback(signDataMap, cli.Params.RSAPrivateKey)
		if err != nil || !verifyResult {
			return fmt.Errorf("verify signature failed")
		}

		if cast.ToString(req.Data.MerchantId) != cli.Params.MerchantId || cast.ToString(req.Data.TerminalId) != cli.Params.TerminalId {
			return fmt.Errorf("illegal merchantId")
		}

		if req.Data.Action != 50 {
			return fmt.Errorf("illegal action")
		}

		//2. 业务侧开始处理
		err = processor(req)
		if err != nil {
			return err
		}

		//3. 最后去confirm
		verifyReq := TCPayVerifyPaymentReq{
			Token: req.Data.Token,
		}
		return cli.VerifyPayment(verifyReq)
	} else {
		//失败
		return fmt.Errorf("%d,%s", req.ResCode, req.Description)
	}
}

// 收到callback后15分钟内必须发送
func (cli *Client) VerifyPayment(req TCPayVerifyPaymentReq) error {

	var signDataMap map[string]interface{}
	mapstructure.Decode(req, &signDataMap)

	signStr, _ := cli.signUtil.GetSign(signDataMap, cli.Params.RSAPrivateKey, 2) //私钥加密
	signDataMap["SignData"] = signStr

	//---------------------------------

	var result TCPayVerifyPaymentResp

	resp, err := cli.ryClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetCloseConnection(true).
		R().
		SetBody(signDataMap).
		SetHeaders(getHeaders()).
		SetDebug(cli.debugMode).
		SetResult(&result).
		SetError(&result).
		Post(cli.Params.VerifyPaymentUrl)

	restLog, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(utils.GetRestyLog(resp))
	cli.logger.Infof("PSPResty#tcpay#verify->%s", string(restLog))

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		//反序列化错误会在此捕捉
		return fmt.Errorf("status code: %d", resp.StatusCode())
	}

	if resp.Error() != nil {
		//反序列化错误会在此捕捉
		return fmt.Errorf("%v, body:%s", resp.Error(), resp.Body())
	}

	return nil
}
