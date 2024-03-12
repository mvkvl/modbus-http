package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Register struct {
	Device  *Device `json:"-"`
	Type    RegType `json:"type,string,omitempty"`
	Mode    RegMode `json:"mode,string,omitempty"`
	Title   string  `json:"title,omitempty"`
	Address uint16  `json:"address,omitempty"`
	Size    uint16  `json:"size,omitempty"`
	Factor  float32 `json:"factor,omitempty"`
}

func (r Register) String() string {
	return fmt.Sprintf("type: %s, mode: %s, addr: %d, size: %d", r.Type, r.Mode, r.Address, r.Size)
}

// UnmarshalJSON custom deserializer to apply default values in case of empty fields
func (r *Register) UnmarshalJSON(data []byte) (err error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	var register Register

	if nil != obj["type"] {
		register.Type, _ = parseRegType(fmt.Sprint(obj["type"]))
	}
	if nil != obj["mode"] {
		register.Mode, _ = parseRegMode(fmt.Sprint(obj["mode"]))
	} else if nil != obj["type"] {
		if register.Type == COIL || register.Type == HOLDING {
			register.Mode = RW
		} else {
			register.Mode = RO
		}
	}
	if nil != obj["address"] {
		v, _ := strconv.Atoi(fmt.Sprint(obj["address"]))
		register.Address = uint16(v)
	}
	if nil != obj["title"] {
		register.Title = fmt.Sprint(obj["title"])
	}
	if nil != obj["size"] {
		v, _ := strconv.Atoi(fmt.Sprint(obj["size"]))
		register.Size = uint16(v)
	} else {
		register.Size = 1
	}
	if nil != obj["factor"] {
		v, _ := strconv.ParseFloat(fmt.Sprint(obj["factor"]), 32)
		if err != nil {
			//log.Fatalf("%q\n", err)
			v = 1.0
		}
		register.Factor = float32(v)
	} else {
		register.Factor = 1.0
	}
	*r = register
	return nil
}
