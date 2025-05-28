package utils

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"reflect"
	"sort"
	"strings"
)

type CheezeebitRSASignatureUtil struct{}

func (util *CheezeebitRSASignatureUtil) GetSign(paramsMap map[string]interface{}, privateKey string) (string, error) {
	delete(paramsMap, "platSign")
	textContent := util.GetContent(paramsMap)
	return util.Sign(textContent, privateKey)
}

func (util *CheezeebitRSASignatureUtil) VerifySign(paramsMap map[string]interface{}, publicKey string, sign string) (bool, error) {
	delete(paramsMap, "platSign")
	textContent := util.GetVerifyContent(paramsMap)
	return util.Verify(textContent, sign, publicKey)
}

//-------------------------------------------------------------

func (util *CheezeebitRSASignatureUtil) GetContent(paramsMap map[string]interface{}) string {
	// Get sorted keys
	keys := lo.Keys(paramsMap)
	sort.Strings(keys)

	var pairs []string
	lo.ForEach(keys, func(x string, index int) {
		value := ""
		if x != "payeeAccountInfos" {
			if x == "agentOrderBatch" {
				//官方文档中把这句给注释掉了
				//valueByte, _ := json.Marshal(paramsMap[x])
				//value = string(valueByte)
			} else {
				value = cast.ToString(paramsMap[x])
			}
		}

		if value != "" {
			pairs = append(pairs, value)
		}
	})

	queryString := strings.Join(pairs, "")
	fmt.Printf("[rawString]%s\n", queryString)

	return queryString
}

func (util *CheezeebitRSASignatureUtil) GetVerifyContent(params map[string]interface{}) string {
	// 获取并排序键名
	paramNames := lo.Keys(params)
	sort.Strings(paramNames)

	var builder strings.Builder

	for _, name := range paramNames {
		value := params[name]
		if !isEmpty(value) {
			switch name {
			case "data":
				if isMap(value) {
					// 递归处理嵌套的 map
					if nestedMap, ok := value.(map[string]interface{}); ok {
						builder.WriteString(util.GetVerifyContent(nestedMap))
					}
				} else {
					builder.WriteString(cast.ToString(value))
				}
			case "payeeAccountInfos":
				// 只有 string 类型才参与签名
				if _, ok := value.(string); ok {
					builder.WriteString(cast.ToString(value))
				}
			default:
				if value != nil {
					builder.WriteString(cast.ToString(value))
				}
			}
		}
	}

	return builder.String()
}

// 辅助函数：检查值是否为空
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return v.Len() == 0
	default:
		return false
	}
}

// 辅助函数：检查是否是 map 类型
func isMap(value interface{}) bool {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Map
}

//---------

func (util *CheezeebitRSASignatureUtil) Sign(message, privateKeyString string) (string, error) {

	signResult, err := SignSHA256RSA([]byte(message), privateKeyString)
	if err != nil {
		fmt.Printf("==sign===>%s\n", err.Error())
	}
	return signResult, nil
}

func (util *CheezeebitRSASignatureUtil) Verify(message, signatureString, publicKeyString string) (bool, error) {

	return VerifySHA256RSA([]byte(message), publicKeyString, signatureString)

}
