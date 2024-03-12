package model

import (
	"encoding/json"
	"testing"
)

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
