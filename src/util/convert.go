package util

import (
	"errors"
	"math"
	"strings"
)

func ToFloat64(unk any) (float64, error) {
	switch i := unk.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int16:
		return float64(i), nil
	case int8:
		return float64(i), nil
	default:
		return math.NaN(), errors.New("convert: unknown value is of incompatible type")
	}
}

func HexaNumberToInteger(hexaString string) (string, error) {
	// replace 0x or 0X with empty String
	if strings.Contains(hexaString, "0x") && 0 == strings.Index(hexaString, "0x") ||
		strings.Contains(hexaString, "0X") && 0 == strings.Index(hexaString, "0X") {
		numberStr := strings.Replace(hexaString, "0x", "", -1)
		numberStr = strings.Replace(numberStr, "0X", "", -1)
		return numberStr, nil
	}
	return "", errors.New("not a hexadecimal string")
}
