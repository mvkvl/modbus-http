package device

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/maja42/goval"
	"github.com/mvkvl/modbus"
	"github.com/mvkvl/modbus-http/model"
	"math"
	"strings"
)

// region - public API

func getReader(register *model.Register, client *modbus.Client) func(slaveId uint8, address, quantity uint16) (results []byte, err error) {
	switch register.Type {
	case model.COIL:
		return (*client).ReadCoils
	case model.DISCRETE:
		return (*client).ReadDiscreteInputs
	case model.INPUT:
		return (*client).ReadInputRegisters
	default:
		return (*client).ReadHoldingRegisters
	}
}

func findChannelByTitle(config *model.Config, title string) (*model.Channel, error) {
	for _, v := range config.Channels {
		if v.Title == title {
			return &v, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no channel found for title '%s'", title))
}
func findDeviceByTitle(channel *model.Channel, title string) (*model.Device, error) {
	for _, v := range channel.Devices {
		if v.Title == title {
			return &v, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no device found for title '%s'", title))
}
func findRegisterByTitle(device *model.Device, title string) (*model.Register, error) {
	for _, v := range device.Registers {
		if v.Title == title {
			return &v, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no register found for title '%s'", title))
}
func findRegister(config *model.Config, reference string) (*model.Register, error) {
	ref := strings.Split(reference, ":")
	if 3 != len(ref) {
		return nil, errors.New(fmt.Sprintf("invalid reference passed: '%s'", reference))
	}
	c, err := findChannelByTitle(config, strings.TrimSpace(ref[0]))
	if nil != err {
		return nil, err
	}
	d, err := findDeviceByTitle(c, strings.TrimSpace(ref[1]))
	if nil != err {
		return nil, err
	}
	r, err := findRegisterByTitle(d, strings.TrimSpace(ref[2]))
	return r, err
}

func ReadFloatRegister(client *modbus.Client, register *model.Register) (result float64, err error) {
	reader := getReader(register, client)
	buff, err := reader(register.Device.SlaveId, register.Address, register.Size)

	if nil != err {
		return 0, err
	}
	var val uint16
	if 0 == len(buff) {
		return 0, errors.New("no value")
	} else if 1 == len(buff) {
		val = uint16(buff[0])
	} else if 2 == len(buff) {
		val = binary.BigEndian.Uint16(buff)
	} else {
		val = binary.BigEndian.Uint16(buff)
	}
	expression := fmt.Sprintf("%f * %d", register.Factor, val)
	eval := goval.NewEvaluator()
	v, err := eval.Evaluate(expression, nil, nil)

	switch i := v.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	default:
		return math.NaN(), errors.New("readFloat: unknown value is of incompatible type")
	}

}
func ReadFloat(client *modbus.Client, config *model.Config, reference string) (result float64, title string, err error) {
	reg, err := findRegister(config, reference)
	if nil != err {
		return 0, "", err
	}
	v, e := ReadFloatRegister(client, reg)
	return v, reg.Title, e
}

// endregion
