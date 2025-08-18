package go_tcpay

import (
	"fmt"
	"testing"
)

func TestAccountInquiry(t *testing.T) {
	vLog := VLog{}
	//构造client
	cli := NewClient(vLog, &TCPayInitParams{MERCHANT_ID, TERMINAL_ID, RSA_PUBLIC_KEY, RSA_PRIVATE_KEY, GARTWAY_URL, CREATE_PAYMENT_URL, WITHDRAW_URL, VERIFY_PAYMENT_URL, DEPOSIT_CALLBACK_URL, WITHDRAW_CALLBACK_URL})

	//发请求
	cli.SetDebugModel(true)
	resp, err := cli.AccountInquiry(GenAccountInquiryRequestDemo())
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
	fmt.Printf("resp:%+v\n", resp)
}

func GenAccountInquiryRequestDemo() TCPayAccountInquiryReq {
	return TCPayAccountInquiryReq{
		AccountNumber: "USD1071343",
	}
}

func TestAutoWithdraw(t *testing.T) {

	vLog := VLog{}
	//构造client
	cli := NewClient(vLog, &TCPayInitParams{MERCHANT_ID, TERMINAL_ID, RSA_PUBLIC_KEY, RSA_PRIVATE_KEY, GARTWAY_URL, CREATE_PAYMENT_URL, WITHDRAW_URL, VERIFY_PAYMENT_URL, DEPOSIT_CALLBACK_URL, WITHDRAW_CALLBACK_URL})

	//发请求
	cli.SetDebugModel(true)
	resp, err := cli.AutoWithdraw(GenWithdrawRequestDemo())
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
	fmt.Printf("resp:%+v\n", resp)
}

func GenWithdrawRequestDemo() TCPayAutoWithdrawPaymentReq {
	return TCPayAutoWithdrawPaymentReq{
		PaymentNumber: "202508180913120744",
		Amount:        "20",         //不能浮点数
		AccountNumber: "USD1071343", //"amirdanialsaedi@gmail.com",
	}
}
