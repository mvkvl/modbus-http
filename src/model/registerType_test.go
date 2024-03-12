package model

import (
	"encoding/json"
	"testing"
)

func TestRegTypeValueSerialization(t *testing.T) {
	var rt = COIL
	exp := "\"coil\""
	data, _ := json.Marshal(rt)
	res := string(data)
	if res != exp {
		t.Errorf("expected '%s', got '%s' instead", exp, res)
	}
}
func TestRegTypeValueDeserialization(t *testing.T) {
	exp := COIL
	data := []byte("{\"type\": \"coil\"}")
	var register Register
	if err := json.Unmarshal(data, &register); err != nil {
		t.Fatalf("%s", err)
	}
	if exp != register.Type {
		t.Errorf("expected type %d (%s), got %d instead", exp, regTypeName[uint8(exp)], register.Type)
	}
}
