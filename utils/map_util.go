package utils

import "github.com/spf13/cast"

func ConvertToStringMap(input map[string]interface{}) map[string]string {
	output := make(map[string]string)
	for key, value := range input {
		output[key] = cast.ToString(value) //fmt.Sprintf("%v", value)
	}
	return output
}
