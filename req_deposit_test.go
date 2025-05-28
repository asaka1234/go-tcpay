package go_cheezeepay

import (
	"fmt"
	"testing"
)

func TestDeposit(t *testing.T) {

	//构造client
	cli := NewClient(nil, MERCHANT_ID, RSA_PUBLIC_KEY, RSA_PRIVATE_KEY, DEPOST_URL, WITHDRAW_URL, DEPOST_BACK_URL, WITHDRAW_BACK_URL)

	//发请求
	resp, err := cli.Deposit(GenDepositRequestDemo())
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
	fmt.Printf("resp:%+v\n", resp)
}

func GenDepositRequestDemo() CheezeePayDepositReq {
	return CheezeePayDepositReq{
		CustomerMerchantsId: "12345", //商户uid
		LegalCoin:           "INR",
		MerchantOrderId:     "8787791",
		DealAmount:          "200", //不能浮点数
	}
}
