package go_tcpay

import (
	"github.com/asaka1234/go-tcpay/utils"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	MerchantID    string // merchantId
	TerminalID    string
	RSAPublicKey  string // 公钥(貌似没用到)
	RSAPrivateKey string // 私钥

	GatewayURL          string //带token跳转地址
	CreatePaymentURL    string
	VerifyPaymentURL    string
	DepositCallbackURL  string
	WithdrawCallbackURL string

	ryClient *resty.Client
	logger   utils.Logger
	signUtil utils.TCPayRSASignatureUtil
}

func NewClient(logger utils.Logger, merchantID, terminalID string, rsaPublicKey, rsaPrivateKey, gatewayURL, createPaymentURL, verifyPaymentURL, depositCallbackURL, withdrawCallbackURL string) *Client {
	return &Client{
		MerchantID:          merchantID,
		TerminalID:          terminalID,
		RSAPublicKey:        rsaPublicKey,
		RSAPrivateKey:       rsaPrivateKey,
		GatewayURL:          gatewayURL,
		CreatePaymentURL:    createPaymentURL,
		VerifyPaymentURL:    verifyPaymentURL,
		DepositCallbackURL:  depositCallbackURL,
		WithdrawCallbackURL: withdrawCallbackURL,

		ryClient: resty.New(), //client实例
		logger:   logger,
		signUtil: utils.TCPayRSASignatureUtil{},
	}
}
