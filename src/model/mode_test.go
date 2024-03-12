package model

import (
	"encoding/json"
	"testing"
)

func TestModeValueSerialization(t *testing.T) {
	var mode = RTU
	exp := "\"rtu\""
	data, _ := json.Marshal(mode)
	res := string(data)
	if res != exp {
		t.Errorf("expected '%s', got '%s' instead", exp, res)
	}
}
func TestModeValueDeserialization(t *testing.T) {
	exp := RTU
	data := []byte("{\"mode\": \"rtu\", \"connection\": \"localhost:20108\"}")
	var channel Channel
	if err := json.Unmarshal(data, &channel); err != nil {
		t.Fatalf("%s", err)
	}
	if exp != channel.Mode {
		t.Errorf("expected mode %d (%s), got %d instead", exp, modeName[uint8(exp)], channel.Mode)
	}
}
