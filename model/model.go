package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// region => enums

// region - Mode enum

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

// endregion

// region - Register Type enum
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

// endregion

// region - Register Mode enum
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

// endregion

// endregion

// region => config elements

// region - Config

type Config struct {
	Channels []Channel `json:"channels,omitempty"`
}

// endregion

// region - Channel

const (
	DefaultCyclePollPause    = 100
	DefaultRegisterPollPause = 10
)

type Channel struct {
	Mode          Mode     `json:"mode,omitempty"`
	Title         string   `json:"title,omitempty"`
	Connection    string   `json:"connection,omitempty"`
	CyclePause    int      `json:"cycle_pause,omitempty"`
	RegisterPause int      `json:"register_pause,omitempty"`
	Devices       []Device `json:"devices,omitempty"`
}

func (c Channel) String() string {
	return fmt.Sprintf(
		"mode: %s, conn: %s, devices: %d, cpause: %d, rpause: %d",
		c.Mode, c.Connection, len(c.Devices), c.GetCyclePause(), c.GetRegisterPause(),
	)
}

func (c Channel) GetCyclePause() int {
	if c.CyclePause <= 0 {
		return DefaultCyclePollPause
	} else {
		return c.CyclePause
	}
}

func (c Channel) GetRegisterPause() int {
	if c.RegisterPause <= 0 {
		return DefaultRegisterPollPause
	} else {
		return c.RegisterPause
	}
}

// endregion

// region - Device

type Device struct {
	SlaveId   uint8      `json:"slave_id,omitempty"`
	Title     string     `json:"title,omitempty"`
	Registers []Register `json:"registers,omitempty"`
}

// endregion

// region - Register

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

// endregion

// endregion
