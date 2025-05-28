package go_tcpay

import (
	"fmt"
	"testing"
)

func TestDeposit(t *testing.T) {

	//构造client
	cli := NewClient(nil, MERCHANT_ID, TERMINAL_ID, RSA_PRIVATE_KEY, RSA_PRIVATE_KEY, GARTWAY_URL, CREATE_PAYMENT_URL, VERIFY_PAYMENT_URL, DEPOSIT_CALLBACK_URL, WITHDRAW_CALLBACK_URL)

	//发请求
	resp, err := cli.Deposit(GenDepositRequestDemo())
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
	fmt.Printf("resp:%+v\n", resp)
}

func GenDepositRequestDemo() TCPayCreatePaymentReq {
	return TCPayCreatePaymentReq{
		ConsumerId:    "12345", //商户uid
		InvoiceNumber: 8787791,
		Amount:        200.00, //不能浮点数
	}
}
