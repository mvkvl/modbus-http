package model

import (
	"encoding/json"
	"testing"
)

func TestRegModeValueSerialization(t *testing.T) {
	var rm = RW
	exp := "\"rw\""
	data, _ := json.Marshal(rm)
	res := string(data)
	if res != exp {
		t.Errorf("expected '%s', got '%s' instead", exp, res)
	}
}
func TestRegModeValueDeserialization(t *testing.T) {
	exp := RW
	data := []byte("{\"mode\": \"rw\"}")
	var register Register
	if err := json.Unmarshal(data, &register); err != nil {
		t.Fatalf("%s", err)
	}
	if exp != register.Mode {
		t.Errorf("expected type %d (%s), got %d instead", exp, regModeName[uint8(exp)], register.Type)
	}
}
