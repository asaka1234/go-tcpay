package go_tcpay

import (
	"fmt"
	"testing"
)

func TestDepositCallback(t *testing.T) {

	vLog := VLog{}
	//构造client
	cli := NewClient(vLog, &TCPayInitParams{MERCHANT_ID, TERMINAL_ID, RSA_PRIVATE_KEY, RSA_PRIVATE_KEY, GARTWAY_URL, CREATE_PAYMENT_URL, WITHDRAW_URL, VERIFY_PAYMENT_URL, DEPOSIT_CALLBACK_URL, WITHDRAW_CALLBACK_URL})

	//发请求
	err := cli.DepositCallback(GenDepositCallbackRequestDemo(), processor)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
}

// Origin=Website&resCode=0&description=Success&data.Action=50
// data.AdditionalData=8df3d80cabd794c3f658b3afe0dd934a&data.Amount=57.00&data.InvoiceNumber=202507230915400460
// data.MerchantId=200157&data.TerminalId=300384&data.Token=ui1nwl4uxn14STbXF71Pbe2RmrXPedR2PxfrVWh4kZc

//{"MerchantId":"200157","TerminalId":"300384","Amount":"24.00","Action":"50","InvoiceNumber":"202507211626030579","TransactionId":"","token":"","AdditionalData":"25f5bed15931dbe97776edf4211f263a"}

func GenDepositCallbackRequestDemo() TCPayCreatePaymentBackReq {
	return TCPayCreatePaymentBackReq{
		ResCode:     0, //商户uid
		Description: "Success",
		Data: &TCPayCreatePaymentBackReqData{
			Action:         "50",
			AdditionalData: "bc1a0e9f324c34db4b88f3d8027a1e08",
			Amount:         "201.00",
			InvoiceNumber:  "202507242147030906",
			MerchantId:     "200157",
			TerminalId:     "300384",
			Token:          "YmC4cfhd2oCRlqXr3Upfd2qgpD0aqepzSimCIh_FJFA",
		}, //不能浮点数
	}
}

//Origin=Website&resCode=0&description=Success&data.Action=50&data.AdditionalData=bc1a0e9f324c34db4b88f3d8027a1e08&data.Amount=201.00&data.InvoiceNumber=202507242147030906&data.MerchantId=200157&data.TerminalId=300384&data.Token=YmC4cfhd2oCRlqXr3Upfd2qgpD0aqepzSimCIh_FJFA

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
