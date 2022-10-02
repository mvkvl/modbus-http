package service

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

// region - reader type

type Reader struct {
	Config *model.Config
	Client *modbus.Client
}

func NewReader(config *model.Config, client *modbus.Client) ReaderAPI {
	return &Reader{
		Config: config,
		Client: client,
	}
}

// endregion

// region - public API

type ReaderAPI interface {
	ReadFloatRegister(register *model.Register) (result float64, err error)
	ReadFloat(reference string) (result float64, title string, err error)
}

func (reader *Reader) ReadFloatRegister(register *model.Register) (result float64, err error) {
	buff, err := reader.getReaderFunction(register)(register.Device.SlaveId, register.Address, register.Size)
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
func (reader *Reader) ReadFloat(reference string) (result float64, title string, err error) {
	reg, err := reader.findRegister(reference)
	if nil != err {
		return 0, "", err
	}
	v, e := reader.ReadFloatRegister(reg)
	return v, reg.Title, e
}

// endregion

// region - private methods

func (reader *Reader) getReaderFunction(register *model.Register) func(slaveId uint8, address, quantity uint16) (results []byte, err error) {
	switch register.Type {
	case model.COIL:
		return (*reader.Client).ReadCoils
	case model.DISCRETE:
		return (*reader.Client).ReadDiscreteInputs
	case model.INPUT:
		return (*reader.Client).ReadInputRegisters
	default:
		return (*reader.Client).ReadHoldingRegisters
	}
}

func (reader *Reader) findChannelByTitle(title string) (*model.Channel, error) {
	for _, v := range reader.Config.Channels {
		if v.Title == title {
			return &v, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no channel found for title '%s'", title))
}
func (reader *Reader) findDeviceByTitle(channel *model.Channel, title string) (*model.Device, error) {
	for _, v := range channel.Devices {
		if v.Title == title {
			return &v, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no device found for title '%s'", title))
}
func (reader *Reader) findRegisterByTitle(device *model.Device, title string) (*model.Register, error) {
	for _, v := range device.Registers {
		if v.Title == title {
			return &v, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no register found for title '%s'", title))
}
func (reader *Reader) findRegister(reference string) (*model.Register, error) {
	ref := strings.Split(reference, ":")
	if 3 != len(ref) {
		return nil, errors.New(fmt.Sprintf("invalid reference passed: '%s'", reference))
	}
	c, err := reader.findChannelByTitle(strings.TrimSpace(ref[0]))
	if nil != err {
		return nil, err
	}
	d, err := reader.findDeviceByTitle(c, strings.TrimSpace(ref[1]))
	if nil != err {
		return nil, err
	}
	r, err := reader.findRegisterByTitle(d, strings.TrimSpace(ref[2]))
	return r, err
}

// endregion
