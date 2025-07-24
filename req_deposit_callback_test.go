package go_tcpay

import (
	"fmt"
	"testing"
)

func TestDepositCallback(t *testing.T) {

	vLog := VLog{}
	//构造client
	cli := NewClient(vLog, &TCPayInitParams{MERCHANT_ID, TERMINAL_ID, RSA_PRIVATE_KEY, RSA_PRIVATE_KEY, GARTWAY_URL, CREATE_PAYMENT_URL, VERIFY_PAYMENT_URL, DEPOSIT_CALLBACK_URL, WITHDRAW_CALLBACK_URL})

	//发请求
	err := cli.DepositCallback(GenDepositCallbackRequestDemo2(), processor)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
}

// Origin=Website&resCode=0&description=Success&data.Action=50
// data.AdditionalData=8df3d80cabd794c3f658b3afe0dd934a&data.Amount=57.00&data.InvoiceNumber=202507230915400460
// data.MerchantId=200157&data.TerminalId=300384&data.Token=ui1nwl4uxn14STbXF71Pbe2RmrXPedR2PxfrVWh4kZc
func GenDepositCallbackRequestDemo() TCPayCreatePaymentBackReq {
	return TCPayCreatePaymentBackReq{
		ResCode:     0, //商户uid
		Description: "Success",
		Data: &TCPayCreatePaymentBackReqData{
			Action:         "50",
			AdditionalData: "8df3d80cabd794c3f658b3afe0dd934a",
			Amount:         "57.00",
			InvoiceNumber:  "202507230915400460",
			MerchantId:     "200157",
			TerminalId:     "300384",
			Token:          "ui1nwl4uxn14STbXF71Pbe2RmrXPedR2PxfrVWh4kZc",
		}, //不能浮点数
	}
}

// Origin=Website&resCode=418&description=The+time+allowed+for+payment+has+expired.&data=
func GenDepositCallbackRequestDemo2() TCPayCreatePaymentBackReq {
	return TCPayCreatePaymentBackReq{
		ResCode:     418, //商户uid
		Description: "hahah",
		Data:        nil, //不能浮点数
	}
}

func processor(TCPayCreatePaymentBackReq) error {
	return nil
}
