package go_tcpay

import "time"

type TCPayInitParams struct {
	MerchantId    string `json:"merchantId" mapstructure:"merchantId" config:"merchantId" yaml:"merchantId"` // merchantId
	TerminalId    string `json:"terminalId" mapstructure:"terminalId" config:"terminalId" yaml:"terminalId"`
	RSAPublicKey  string `json:"rsaPublicKey" mapstructure:"rsaPublicKey" config:"rsaPublicKey" yaml:"rsaPublicKey"`     // 公钥(貌似没用到)
	RSAPrivateKey string `json:"rsaPrivateKey" mapstructure:"rsaPrivateKey" config:"rsaPrivateKey" yaml:"rsaPrivateKey"` // 私钥

	GatewayUrl       string `json:"gatewayUrl" mapstructure:"gatewayUrl" config:"gatewayUrl" yaml:"gatewayUrl"` //带token跳转地址
	CreatePaymentUrl string `json:"createPaymentUrl" mapstructure:"createPaymentUrl" config:"createPaymentUrl" yaml:"createPaymentUrl"`
	WithdrawUrl      string `json:"withdrawUrl" mapstructure:"withdrawUrl" config:"withdrawUrl" yaml:"withdrawUrl"`

	VerifyPaymentUrl string `json:"verifyPaymentUrl" mapstructure:"verifyPaymentUrl" config:"verifyPaymentUrl" yaml:"verifyPaymentUrl"`
	DepositBackUrl   string `json:"depositBackUrl" mapstructure:"depositBackUrl" config:"depositBackUrl" yaml:"depositBackUrl"`
	WithdrawBackUrl  string `json:"withdrawBackUrl" mapstructure:"withdrawBackUrl" config:"withdrawBackUrl" yaml:"withdrawBackUrl"`
}

// ---------------------------------------------
// create payment
type TCPayCreatePaymentReq struct {
	ConsumerId    string `json:"ConsumerId,omitempty" mapstructure:"ConsumerId"`
	Amount        string `json:"Amount" mapstructure:"Amount"`               // 要求:必须2位精度!! 哪怕是整数,也要1.00这样
	InvoiceNumber int64  `json:"InvoiceNumber" mapstructure:"InvoiceNumber"` //商户订单号
	//以下sdk来设置
	//MerchantId    int    `json:"MerchantId" mapstructure:"MerchantId"`       //商户号
	//TerminalId    int    `json:"TerminalId" mapstructure:"TerminalId"`       //终端号
	//LocalDateTime string `json:"LocalDateTime" mapstructure:"LocalDateTime"` //yyyy/MM/dd HH:mm:ss
	//AdditionalData string  `json:"additionalData,omitempty" mapstructure:"additionalData"` //option  备注信息,用来做一个简单签名,解决callback没有鉴权的问题. Additional transaction information
	//Action        int     `json:"Action" mapstructure:"Action"` // 50-deposit, 100-withdrawl
	//ReturnUrl     string  `json:"ReturnUrl" mapstructure:"ReturnUrl"`         //是ajax callback地址
	//SignData      string `json:"SignData" mapstructure:"SignData"`           //签名数据
}

// auto withdraw
type TCPayAutoWithdrawPaymentReq struct {
	Amount        string `json:"Amount" mapstructure:"Amount"`               // 要求:必须2位精度!! 哪怕是整数,也要1.00这样
	AccountNumber string `json:"AccountNumber" mapstructure:"AccountNumber"` //客户在TcPay的用户名
	PaymentNumber string `json:"PaymentNumber" mapstructure:"PaymentNumber"` //商户订单号
	//以下sdk来设置
	//MerchantId    int    `json:"MerchantId" mapstructure:"MerchantId"`       //商户号
	//TerminalId    int    `json:"TerminalId" mapstructure:"TerminalId"`       //终端号
	//AdditionalData string  `json:"additionalData,omitempty" mapstructure:"additionalData"` //option  备注信息,用来做一个简单签名,解决callback没有鉴权的问题. Additional transaction information
	//SignData      string `json:"SignData" mapstructure:"SignData"`           //签名数据
}

// auto withdraw
type TCPayAutoWithdrawDetailReq struct {
	TransactionId int `json:"TransactionId" mapstructure:"TransactionId"` // 交易ID
	//以下sdk来设置
	//MerchantId    int    `json:"MerchantId" mapstructure:"MerchantId"`       //商户号
	//TerminalId    int    `json:"TerminalId" mapstructure:"TerminalId"`       //终端号
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
// multipart/form-data  或者 x-www-form-urlencoded , 所以必须要指定form标签
// 其中data里的嵌套字段,指定时要格外注意
type TCPayCreatePaymentBackReq struct {
	ResCode     int                            `json:"resCode" form:"resCode"  mapstructure:"resCode"` //0是成功
	Description string                         `json:"description" form:"description" mapstructure:"description"`
	Data        *TCPayCreatePaymentBackReqData `json:"data,omitempty" form:"data" mapstructure:"data"`
}

type TCPayCreatePaymentCommonBackReq struct {
	ResCode     int    `json:"resCode" form:"resCode"  mapstructure:"resCode"` //0是成功
	Description string `json:"description" form:"description" mapstructure:"description"`
}

type TCPayCreatePaymentBackReqData struct {
	MerchantId     string `json:"MerchantId" form:"data.MerchantId" mapstructure:"MerchantId"` //嵌套
	TerminalId     string `json:"TerminalId" form:"data.TerminalId" mapstructure:"TerminalId"`
	Amount         string `json:"Amount" form:"data.Amount" mapstructure:"Amount"`
	Action         string `json:"Action" form:"data.Action" mapstructure:"Action"`                      // 50-deposit, 100-withdrawl
	InvoiceNumber  string `json:"InvoiceNumber" form:"data.InvoiceNumber" mapstructure:"InvoiceNumber"` //商户订单号
	TransactionId  string `json:"TransactionId" form:"data.TransactionId" mapstructure:"TransactionId"` //TODO 这个字段文档中没有,要确认下
	Token          string `json:"Token" form:"data.Token" mapstructure:"Token"`
	AdditionalData string `json:"AdditionalData,omitempty" form:"data.AdditionalData" mapstructure:"AdditionalData"` //我觉得这里最好还是要用来做一下sign签名,不然还是很容易被伪造得
}

//=====================最终确认=========================

type TCPayVerifyPaymentReq struct {
	Token string `json:"Token" mapstructure:"Token"`
	//以下是sdk来计算
	//SignData string `json:"SignData" mapstructure:"SignData"`
}

type TCPayVerifyPaymentResp struct {
	ResCode     int         `json:"ResCode" mapstructure:"ResCode"` //0是成功
	Description string      `json:"Description" mapstructure:"Description"`
	Data        interface{} `json:"Data,omitempty" mapstructure:"Data"`
}

// 错误的时候:是slice.  正确的时候是如下信息(但是这个data数据无关紧要,所以先不处理了)
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
