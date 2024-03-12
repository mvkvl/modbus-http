package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Mode uint8

const (
	RTU Mode = iota + 1
	TCP
	ENC
)

var (
	modeName = map[uint8]string{
		1: "rtu",
		2: "tcp",
		3: "enc",
	}
	modeValue = map[string]uint8{
		"rtu": 1,
		"tcp": 2,
		"enc": 3,
	}
)

func parseMode(s string) (Mode, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := modeValue[s]
	if !ok {
		return Mode(0), fmt.Errorf("%q is not a valid channel operation mode", s)
	}
	return Mode(value), nil
}
func (m Mode) String() string {
	return modeName[uint8(m)]
}
func (m Mode) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}
func (m *Mode) UnmarshalJSON(data []byte) (err error) {
	var input string
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}
	if *m, err = parseMode(input); err != nil {
		return err
	}
	return nil
}
