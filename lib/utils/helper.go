package utils

import "fmt"

func ConvertToString(i interface{}) string {
	if i == nil {
		return ""
	}

	str := fmt.Sprintf("%v", i)
	return str
}
