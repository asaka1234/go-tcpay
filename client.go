package go_tcpay

import (
	"github.com/asaka1234/go-tcpay/utils"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	MerchantID    string // merchantId
	RSAPublicKey  string // 公钥
	RSAPrivateKey string // 私钥

	DepositURL         string
	DepositCallbackURL string

	WithdrawURL         string
	WithdrawCallbackURL string

	ryClient *resty.Client
	logger   utils.Logger
	rsaUtil  utils.CheezeebitRSASignatureUtil
}

func NewClient(logger utils.Logger, merchantID string, rsaPublicKey, rsaPrivateKey, depositURL, withdrawURL, depositCallbackURL, withdrawCallbackURL string) *Client {
	return &Client{
		MerchantID:          merchantID,
		RSAPublicKey:        rsaPublicKey,
		RSAPrivateKey:       rsaPrivateKey,
		DepositURL:          depositURL,
		DepositCallbackURL:  depositCallbackURL,
		WithdrawURL:         withdrawURL,
		WithdrawCallbackURL: withdrawCallbackURL,

		ryClient: resty.New(), //client实例
		logger:   logger,
		rsaUtil:  utils.CheezeebitRSASignatureUtil{},
	}
}
