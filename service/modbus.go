package service

import (
	"encoding/binary"
	"errors"
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/maja42/goval"
	"github.com/mvkvl/modbus"
	"github.com/mvkvl/modbus-http/model"
	"math"
)

// region - reader type

type modbusClient struct {
	Config *model.Config
	Client *modbus.Client
}

type Reader interface {
	Read(register *model.Register) (result float64, err error)
	ReadRef(reference string) (result float64, title string, err error)
}

type Writer interface {
	Write(register *model.Register, value uint16) (err error)
	WriteRef(reference string, value uint16) (err error)
}

type ModbusClient interface {
	Reader
	Writer
}

func NewModbusClient(config *model.Config, client *modbus.Client) ModbusClient {
	return &modbusClient{
		Config: config,
		Client: client,
	}
}

// endregion

// region - public API

// region -> read

func (client *modbusClient) Read(register *model.Register) (result float64, err error) {
	buff, err := client.getReaderFunction(register)(register.Device.SlaveId, register.Address, register.Size)
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
func (client *modbusClient) ReadRef(reference string) (result float64, title string, err error) {
	reg, err := client.Config.FindRegister(reference)
	if nil != err {
		return 0, "", err
	}
	v, e := client.Read(reg)
	return v, reg.Title, e
}

// endregion

// region -> write

func (client *modbusClient) Write(register *model.Register, value uint16) (err error) {
	f := client.getWriterFunction(register)
	if register.Type == model.COIL && value > 0 {
		value = 0xFF00
	}
	buff, err := f(register.Device.SlaveId, register.Address, value)
	if nil == err {
		log.Info("write response: % x\n", buff)
	}
	return err
}
func (client *modbusClient) WriteRef(reference string, value uint16) (err error) {
	reg, err := client.Config.FindRegister(reference)
	if nil != err {
		return err
	}
	if reg.Mode == model.RO {
		return errors.New("trying to write to read only register")
	}
	return client.Write(reg, value)
}

// endregion

// endregion

// region - private methods

func (client *modbusClient) getReaderFunction(register *model.Register) func(slaveId uint8, address, quantity uint16) (results []byte, err error) {
	switch register.Type {
	case model.COIL:
		return (*client.Client).ReadCoils
	case model.DISCRETE:
		return (*client.Client).ReadDiscreteInputs
	case model.INPUT:
		return (*client.Client).ReadInputRegisters
	default:
		return (*client.Client).ReadHoldingRegisters
	}
}
func (client *modbusClient) getWriterFunction(register *model.Register) func(slaveId uint8, address, value uint16) (results []byte, err error) {
	switch register.Type {
	case model.COIL:
		return (*client.Client).WriteSingleCoil
	case model.HOLDING:
		return (*client.Client).WriteSingleRegister
	default:
		panic(fmt.Sprintf("invalid register type used for data write: %s", register.Type))
	}
}

// endregion
