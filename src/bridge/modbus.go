package bridge

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/maja42/goval"
	"github.com/mvkvl/modbus"
	"math"
	"mbridge/model"
	"mbridge/util"
)

// region - API

type Reader interface {
	Read(register *model.Register) (raw uint32, value float64, err error)
	ReadRef(reference string) (raw uint32, value float64, title string, err error)
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
		config: config,
		client: client,
		logger: util.GetLogger("modbus"),
	}
}

// endregion

// region - client implementation

type modbusClient struct {
	config *model.Config
	client *modbus.Client
	logger util.Logger
}

// region => public API
// region ~> read

// Read reads up to 2 words of bus data & converts it to uint32

func (c *modbusClient) Read(register *model.Register) (raw uint32, value float64, err error) {
	reader := c.getReaderFunction(register)
	if reader == nil {
		return 0, 0, fmt.Errorf("nil modbus reader")
	}
	buff, err := reader(register.Device.SlaveId, register.Address, register.Size)
	if nil != err {
		return 0, 0, fmt.Errorf("", err)
	}
	var val uint32
	if 0 == len(buff) {
		return 0, 0, errors.New("no value")
	} else if 1 == len(buff) {
		val = uint32(buff[0])
	} else if 2 == len(buff) {
		val = uint32(binary.BigEndian.Uint16(buff))
	} else if 4 == len(buff) {
		//log.Warn("NxBuff: % 0x, a: % 0x, b: % 0x\n", buff, buff[:2], buff[2:])
		val = uint32(binary.BigEndian.Uint16(buff[:2]))<<8 + uint32(binary.BigEndian.Uint16(buff[2:]))
	} else {
		return 0, math.NaN(), errors.New(fmt.Sprintf("read: too large data chunk received: % 0x", buff))
	}
	expression := fmt.Sprintf("%f * %d", register.Factor, val)
	eval := goval.NewEvaluator()
	v, err := eval.Evaluate(expression, nil, nil)
	switch i := v.(type) {
	case float64:
		return val, i, nil
	case float32:
		return val, float64(i), nil
	case int64:
		return val, float64(i), nil
	default:
		return val, math.NaN(), errors.New("read: unknown value is of incompatible type")
	}
}
func (c *modbusClient) ReadRef(reference string) (raw uint32, value float64, title string, err error) {
	reg, err := c.config.FindRegister(reference)
	if nil != err {
		return 0, 0, "", err
	}
	r, v, e := c.Read(reg)
	return r, v, reg.Title, e
}

// endregion
// region ~> write

func (c *modbusClient) Write(register *model.Register, value uint16) (err error) {
	f := c.getWriterFunction(register)
	if register.Type == model.COIL && value > 0 {
		value = 0xFF00
	}
	buff, err := f(register.Device.SlaveId, register.Address, value)
	if nil == err {
		c.logger.Info("write response: % x", buff)
	}
	return err
}
func (c *modbusClient) WriteRef(reference string, value uint16) (err error) {
	reg, err := c.config.FindRegister(reference)
	if nil != err {
		return err
	}
	if reg.Mode == model.RO {
		return errors.New("trying to write to read only register")
	}
	return c.Write(reg, value)
}

// endregion
// endregion
// region - private methods

func (c *modbusClient) getReaderFunction(register *model.Register) func(slaveId uint8, address, quantity uint16) (results []byte, err error) {
	switch register.Type {
	case model.COIL:
		return (*c.client).ReadCoils
	case model.DISCRETE:
		return (*c.client).ReadDiscreteInputs
	case model.INPUT:
		return (*c.client).ReadInputRegisters
	default:
		return (*c.client).ReadHoldingRegisters
	}
}
func (c *modbusClient) getWriterFunction(register *model.Register) func(slaveId uint8, address, value uint16) (results []byte, err error) {
	switch register.Type {
	case model.COIL:
		return (*c.client).WriteSingleCoil
	case model.HOLDING:
		return (*c.client).WriteSingleRegister
	default:
		panic(fmt.Sprintf("invalid register type used for data write: %s", register.Type))
	}
}

// endregion

// endregion
