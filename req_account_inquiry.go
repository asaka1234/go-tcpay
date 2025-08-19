package go_tcpay

import (
	"crypto/tls"
	"fmt"
	"github.com/asaka1234/go-tcpay/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

func (cli *Client) AutoWithdrawDetail(req TCPayAutoWithdrawDetailReq) (*TCPayCreatePaymentResponse, error) {

	rawURL := "https://pg.toppayment.net/api/v2/Settlement/Details" //cli.Params.AccountInquiryUrl

	// 1. 拿到请求参数，转为map
	var signDataMap map[string]interface{}
	mapstructure.Decode(req, &signDataMap)
	signDataMap["MerchantId"] = cast.ToInt(cli.Params.MerchantId)
	signDataMap["TerminalId"] = cast.ToInt(cli.Params.TerminalId)
	signDataMap["TransactionId"] = req.TransactionId

	// 2. 计算签名,补充参数
	signStr, _ := cli.signUtil.GetSign(signDataMap, cli.Params.RSAPrivateKey, 4) //私钥加密
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
	cli.logger.Infof("PSPResty#tcpay#autoWithdrawDetail->%s", string(restLog))

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
