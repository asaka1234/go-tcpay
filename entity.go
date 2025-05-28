package go_tcpay

import "time"

// ---------------------------------------------
// create payment
type TCPayCreatePaymentReq struct {
	ConsumerId    string  `json:"ConsumerId,omitempty" mapstructure:"ConsumerId"`
	Amount        float64 `json:"Amount" mapstructure:"Amount"`               // 要求:必须2位精度!! 哪怕是整数,也要1.00这样
	InvoiceNumber int64   `json:"InvoiceNumber" mapstructure:"InvoiceNumber"` //商户订单号
	//以下sdk来设置
	//MerchantId    int    `json:"MerchantId" mapstructure:"MerchantId"`       //商户号
	//TerminalId    int    `json:"TerminalId" mapstructure:"TerminalId"`       //终端号
	//LocalDateTime string `json:"LocalDateTime" mapstructure:"LocalDateTime"` //yyyy/MM/dd HH:mm:ss
	//AdditionalData string  `json:"additionalData,omitempty" mapstructure:"additionalData"` //option  备注信息,用来做一个简单签名,解决callback没有鉴权的问题. Additional transaction information
	//Action        int     `json:"Action" mapstructure:"Action"` // 50-deposit, 100-withdrawl
	//ReturnUrl     string  `json:"ReturnUrl" mapstructure:"ReturnUrl"`         //是ajax callback地址
	//SignData      string `json:"SignData" mapstructure:"SignData"`           //签名数据
}

//-------------------------------

type TCPayCreatePaymentResponse struct {
	ResCode     int                             `json:"ResCode" mapstructure:"ResCode"` //0是成功, >0是失败
	Description string                          `json:"Description" mapstructure:"Description"`
	Data        *TCPayCreatePaymentResponseData `json:"Data,omitempty" mapstructure:"Data"`
}

type TCPayCreatePaymentResponseData struct {
	Token string `json:"Token" mapstructure:"Token"`
}

//--------------callback------------------------------

// POST
type TCPayCreatePaymentBackReq struct {
	ResCode     int                            `json:"ResCode" mapstructure:"ResCode"` //0是成功
	Description string                         `json:"Description" mapstructure:"Description"`
	Data        *TCPayCreatePaymentBackReqData `json:"Data,omitempty" mapstructure:"Data"`
}

type TCPayCreatePaymentBackReqData struct {
	MerchantId     int64   `json:"MerchantId" mapstructure:"MerchantId"`
	TerminalId     int64   `json:"TerminalId" mapstructure:"TerminalId"`
	Amount         float64 `json:"Amount" mapstructure:"Amount"`
	Action         int64   `json:"Action" mapstructure:"Action"`               // 50-deposit, 100-withdrawl
	InvoiceNumber  int64   `json:"InvoiceNumber" mapstructure:"InvoiceNumber"` //商户订单号
	TransactionId  int64   `json:"TransactionId" mapstructure:"TransactionId"` //TODO 这个字段文档中没有,要确认下
	Token          string  `json:"token" mapstructure:"token"`
	AdditionalData string  `json:"AdditionalData,omitempty" mapstructure:"AdditionalData"` //我觉得这里最好还是要用来做一下sign签名,不然还是很容易被伪造得
}

//=====================最终确认=========================

type TCPayVerifyPaymentReq struct {
	Token string `json:"Token" mapstructure:"Token"`
	//以下是sdk来计算
	//SignData string `json:"SignData" mapstructure:"SignData"`
}

type TCPayVerifyPaymentResp struct {
	ResCode     int                         `json:"ResCode" mapstructure:"ResCode"` //0是成功
	Description string                      `json:"Description" mapstructure:"Description"`
	Data        *TCPayVerifyPaymentRespData `json:"Data,omitempty" mapstructure:"Data"`
}

type TCPayVerifyPaymentRespData struct {
	MerchantId    int64   `json:"MerchantId" mapstructure:"MerchantId"`
	TerminalId    int64   `json:"TerminalId" mapstructure:"TerminalId"`
	Amount        float64 `json:"Amount" mapstructure:"Amount"`
	Action        int64   `json:"Action" mapstructure:"Action"`
	InvoiceNumber int64   `json:"InvoiceNumber" mapstructure:"InvoiceNumber"`
	TransactionId int64   `json:"TransactionId" mapstructure:"TransactionId"` //是psp的订单号
	//Token          string                      `json:"token" mapstructure:"token"`
	Wallet         string                      `json:"Wallet,omitempty" mapstructure:"Wallet"`
	AdditionalData string                      `json:"AdditionalData,omitempty" mapstructure:"AdditionalData"`
	UserInfo       *TCPayVerifyPaymentUserInfo `json:"UserInfo,omitempty" mapstructure:"UserInfo"`
}

type TCPayVerifyPaymentUserInfo struct {
	FirstName  string    `json:"FirstName" mapstructure:"FirstName"`
	LastName   string    `json:"LastName" mapstructure:"LastName"`
	BirthDate  time.Time `json:"BirthDate" mapstructure:"BirthDate"`
	NationalId string    `json:"NationalId" mapstructure:"NationalId"`
}
