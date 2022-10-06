package service

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Int63()%int64(len(letterRunes))]
	}
	return string(b)
}

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
