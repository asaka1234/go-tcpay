package go_tcpay

import (
	"fmt"
	"github.com/asaka1234/go-tcpay/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

func (cli *Client) WithdrawCallback(req TCPayCreatePaymentBackReq, processor func(TCPayCreatePaymentBackReq) error) error {
	//1. 验证签名
	if req.ResCode == 0 {
		var signDataMap map[string]interface{}
		mapstructure.Decode(req.Data, &signDataMap)

		verifyResult, err := utils.VerifyCallback(signDataMap, cli.RSAPublicKey)
		if err != nil || !verifyResult {
			return fmt.Errorf("illegal callback!")
		}

		if cast.ToString(req.Data.MerchantId) != cli.MerchantID || cast.ToString(req.Data.TerminalId) != cli.TerminalID {
			return fmt.Errorf("illegal merchantID!")
		}

		if req.Data.Action != 100 {
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
