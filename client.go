package go_tcpay

import (
	"github.com/asaka1234/go-tcpay/utils"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	Params *TCPayInitParams

	ryClient  *resty.Client
	logger    utils.Logger
	debugMode bool
	signUtil  utils.TCPayRSASignatureUtil
}

func NewClient(logger utils.Logger, params *TCPayInitParams) *Client {
	return &Client{
		Params: params,

		ryClient:  resty.New(), //client实例
		logger:    logger,
		debugMode: false,
		signUtil:  utils.TCPayRSASignatureUtil{},
	}
}

func (cli *Client) SetDebugModel(debugModel bool) {
	cli.debugMode = debugModel
}
