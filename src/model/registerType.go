package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type RegType uint8

const (
	COIL RegType = iota + 1
	DISCRETE
	INPUT
	HOLDING
)

var (
	regTypeName = map[uint8]string{
		1: "coil",
		2: "discrete",
		3: "input",
		4: "holding",
	}
	regTypeValue = map[string]uint8{
		"coil":     1,
		"discrete": 2,
		"input":    3,
		"holding":  4,
	}
)

func parseRegType(s string) (RegType, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := regTypeValue[s]
	if !ok {
		return RegType(0), fmt.Errorf("%q is not a valid register type", s)
	}
	return RegType(value), nil
}
func (t RegType) String() string {
	return regTypeName[uint8(t)]
}
func (t RegType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
func (t *RegType) UnmarshalJSON(data []byte) (err error) {
	var input string
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}
	if *t, err = parseRegType(input); err != nil {
		return err
	}
	return nil
}
