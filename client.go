package go_tcpay

import (
	"github.com/asaka1234/go-tcpay/utils"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	MerchantID    string // merchantId
	TerminalID    string
	RSAPublicKey  string // 公钥
	RSAPrivateKey string // 私钥

	CreatePaymentURL         string
	VerifyPaymentURL         string
	CreatePaymentCallbackURL string

	ryClient *resty.Client
	logger   utils.Logger
}

func NewClient(logger utils.Logger, merchantID, terminalID string, rsaPublicKey, rsaPrivateKey, createPaymentURL, verifyPaymentURL, createPaymentCallbackURL string) *Client {
	return &Client{
		MerchantID:               merchantID,
		TerminalID:               terminalID,
		RSAPublicKey:             rsaPublicKey,
		RSAPrivateKey:            rsaPrivateKey,
		CreatePaymentURL:         createPaymentURL,
		VerifyPaymentURL:         verifyPaymentURL,
		CreatePaymentCallbackURL: createPaymentCallbackURL,

		ryClient: resty.New(), //client实例
		logger:   logger,
	}
}
