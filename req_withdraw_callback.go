package go_tcpay

import (
	"fmt"
	"github.com/asaka1234/go-tcpay/utils"
	"github.com/spf13/cast"
)

func (cli *Client) WithdrawCallback(req TCPayCreatePaymentBackReq, processor func(TCPayCreatePaymentBackReq) error) error {
	//1. 验证签名
	if req.ResCode == 0 {

		signDataMap := map[string]interface{}{
			"MerchantId":     req.Data.MerchantId,
			"TerminalId":     req.Data.TerminalId,
			"Action":         req.Data.Action,
			"Amount":         req.Data.Amount,
			"InvoiceNumber":  req.Data.InvoiceNumber,
			"AdditionalData": req.Data.AdditionalData,
		}

		verifyResult, err := utils.VerifyCallback(signDataMap, cli.Params.RSAPrivateKey)
		if err != nil || !verifyResult {
			return fmt.Errorf("verify signature failed")
		}

		if cast.ToString(req.Data.MerchantId) != cli.Params.MerchantId || cast.ToString(req.Data.TerminalId) != cli.Params.TerminalId {
			return fmt.Errorf("illegal merchantId")
		}

		if req.Data.Action != "100" {
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
