package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Device struct {
	Channel   *Channel
	SlaveId   uint8      `json:"slave_id,omitempty"`
	Title     string     `json:"title,omitempty"`
	Alias     string     `json:"alias,omitempty"`
	Registers []Register `json:"registers,omitempty"`
}

func (d *Device) findRegisterByTitle(title string) (*Register, error) {
	for _, v := range d.Registers {
		if v.Title == title {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("no register found for title '%s'", title)
}

func (d *Device) UnmarshalJSON(data []byte) (err error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	var device Device

	if nil != obj["title"] {
		device.Title = fmt.Sprint(obj["title"])
	}
	if nil != obj["alias"] {
		device.Alias = fmt.Sprint(obj["alias"])
	} else {
		device.Alias = device.Title
	}
	if nil != obj["slave_id"] {
		v, _ := strconv.Atoi(fmt.Sprint(obj["slave_id"]))
		device.SlaveId = uint8(v)
	}
	if nil != obj["registers"] {
		r := obj["registers"]
		rj, err := json.Marshal(r)
		if nil != err {
			return err
		}
		if err := json.Unmarshal(rj, &device.Registers); err != nil {
			return err
		}
	}
	*d = device
	return nil
}
