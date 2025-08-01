package go_tcpay

import (
	"crypto/tls"
	"fmt"
	"github.com/asaka1234/go-tcpay/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"time"
)

func (cli *Client) Withdraw(req TCPayCreatePaymentReq) (*TCPayCreatePaymentResponse, error) {

	rawURL := cli.Params.CreatePaymentUrl

	// 1. 拿到请求参数，转为map
	var signDataMap map[string]interface{}
	mapstructure.Decode(req, &signDataMap)
	signDataMap["MerchantId"] = cli.Params.MerchantId
	signDataMap["TerminalId"] = cli.Params.TerminalId
	signDataMap["LocalDateTime"] = time.Now().Format("2006/01/02 15:04:05")
	signDataMap["Action"] = 100 //100-withdraw
	signDataMap["ReturnUrl"] = cli.Params.WithdrawBackUrl

	//以此确保amount是2位精度!
	amount := decimal.NewFromFloat(cast.ToFloat64(signDataMap["Amount"])) //转为decimal
	signDataMap["Amount"] = amount.StringFixed(2)

	// 2. 先计算一个md5签名, 随后补充到AdditionalData字段中.
	signDataMap["AdditionalData"], _ = utils.SignCallback(signDataMap, cli.Params.RSAPrivateKey)

	// 2. 计算签名,补充参数
	signStr, _ := cli.signUtil.GetSign(signDataMap, cli.Params.RSAPrivateKey, 1) //私钥加密
	signDataMap["SignData"] = signStr

	var result TCPayCreatePaymentResponse

	resp, err := cli.ryClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetCloseConnection(true).
		R().
		SetBody(signDataMap).
		SetHeaders(getHeaders()).
		SetDebug(cli.debugMode).
		SetResult(&result).
		SetError(&result).
		Post(rawURL)

	restLog, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(utils.GetRestyLog(resp))
	cli.logger.Infof("PSPResty#tcpay#withdraw->%s", string(restLog))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		//反序列化错误会在此捕捉
		return nil, fmt.Errorf("status code: %d", resp.StatusCode())
	}

	if resp.Error() != nil {
		//反序列化错误会在此捕捉
		return nil, fmt.Errorf("%v, body:%s", resp.Error(), resp.Body())
	}

	return &result, nil
}
