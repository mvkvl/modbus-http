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
func TestRegisterDeserializationA(t *testing.T) {
	exp := Register{Mode: RW, Size: 1}
	data := "{\"mode\": \"rw\"}"
	var register Register
	if err := json.Unmarshal([]byte(data), &register); err != nil {
		t.Fatalf("%s", err)
	}
	if exp.Mode != register.Mode {
		t.Errorf("expected mode '%s', got '%s' instead", exp.Mode, register.Mode)
	}
	if exp.Size != register.Size {
		t.Errorf("expected size '%d', got '%d' instead", exp.Size, register.Size)
	}
}
func TestRegisterDeserializationB(t *testing.T) {
	exp := Register{Type: COIL, Mode: RW, Size: 1}
	data := "{\"type\": \"coil\"}"
	var register Register
	if err := json.Unmarshal([]byte(data), &register); err != nil {
		t.Fatalf("%s", err)
	}
	if exp.Mode != register.Mode {
		t.Errorf("expected mode '%s', got '%s' instead\n", exp.Mode, register.Mode)
	}
	if exp.Size != register.Size {
		t.Errorf("expected size '%d', got '%d' instead\n", exp.Size, register.Size)
	}
}
func TestRegisterDeserializationC(t *testing.T) {
	exp := Register{Type: DISCRETE, Mode: RO, Size: 1}
	data := "{\"type\": \"discrete\"}"
	var register Register
	if err := json.Unmarshal([]byte(data), &register); err != nil {
		t.Fatalf("%s", err)
	}
	if exp.Mode != register.Mode {
		t.Errorf("expected mode '%s', got '%s' instead\n", exp.Mode, register.Mode)
	}
	if exp.Size != register.Size {
		t.Errorf("expected size '%d', got '%d' instead\n", exp.Size, register.Size)
	}
}
