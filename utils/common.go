package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func EnsurePrefix(value string, prefix string) string {
	if !strings.HasPrefix(value, prefix) {
		value = fmt.Sprintf("%v%v", prefix, value)
	}
	return value
}

func ParseJSON(value string) (map[string]interface{}, error) {
	var result map[string]interface{}

	err := json.Unmarshal([]byte(value), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}