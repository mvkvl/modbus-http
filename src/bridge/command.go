package bridge

import "mbridge/model"

type Type int

const (
	CTRead Type = iota
	CTWrite
)

type Command interface {
	GetType() Type
	GetChannel() *model.Channel
	GetDevice() *model.Device
	GetRegister() *model.Register
	GetValue() uint16
}

// region - read command

type readCommand struct {
	channel  *model.Channel
	device   *model.Device
	register *model.Register
}

func NewReadCommand(channel *model.Channel, device *model.Device, register *model.Register) Command {
	return &readCommand{
		channel:  channel,
		device:   device,
		register: register,
	}
}

func (c *readCommand) GetType() Type {
	return CTRead
}
func (c *readCommand) GetChannel() *model.Channel {
	return c.channel
}
func (c *readCommand) GetDevice() *model.Device {
	return c.device
}
func (c *readCommand) GetRegister() *model.Register {
	return c.register
}
func (c *readCommand) GetValue() uint16 {
	return 0
}

// endregion
// region - write command

type writeCommand struct {
	channel  *model.Channel
	device   *model.Device
	register *model.Register
	value    uint16
}

func NewWriteCommand(channel *model.Channel, device *model.Device, register *model.Register, value uint16) Command {
	return &writeCommand{
		channel:  channel,
		device:   device,
		register: register,
		value:    value,
	}
}

func (c *writeCommand) GetType() Type {
	return CTWrite
}
func (c *writeCommand) GetChannel() *model.Channel {
	return c.channel
}
func (c *writeCommand) GetDevice() *model.Device {
	return c.device
}
func (c *writeCommand) GetRegister() *model.Register {
	return c.register
}
func (c *writeCommand) GetValue() uint16 {
	return c.value
}

// endregion
