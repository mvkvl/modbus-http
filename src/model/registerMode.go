package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type RegMode uint8

const (
	RO RegMode = iota + 1
	RW
	WO
)

var (
	regModeName = map[uint8]string{
		1: "ro",
		2: "rw",
		3: "wo",
	}
	regModeValue = map[string]uint8{
		"ro": 1,
		"rw": 2,
		"wo": 3,
	}
)

func parseRegMode(s string) (RegMode, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := regModeValue[s]
	if !ok {
		return RegMode(0), fmt.Errorf("%q is not a valid register mode", s)
	}
	return RegMode(value), nil
}
func (t RegMode) String() string {
	return regModeName[uint8(t)]
}
func (t RegMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
func (t *RegMode) UnmarshalJSON(data []byte) (err error) {
	var input string
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}
	if *t, err = parseRegMode(input); err != nil {
		return err
	}
	return nil
}
