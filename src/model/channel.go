package model

import (
	"fmt"
	"time"
)

const (
	defaultCyclePollPause    = time.Millisecond * 100
	defaultRegisterPollPause = time.Millisecond * 10
)

type Channel struct {
	Mode          Mode     `json:"mode,omitempty"`
	Title         string   `json:"title,omitempty"`
	Connection    string   `json:"connection,omitempty"`
	CyclePause    *string  `json:"cycle_pause,omitempty"`
	RegisterPause *string  `json:"register_pause,omitempty"`
	Devices       []Device `json:"devices,omitempty"`
}

func (c Channel) String() string {
	return fmt.Sprintf(
		"mode: %s, conn: %s, devices: %d, cpause: %d, rpause: %d",
		c.Mode, c.Connection, len(c.Devices), c.GetCyclePause(), c.GetRegisterPause(),
	)
}

func (c Channel) GetCyclePause() time.Duration {
	return durationOrDefault(c.CyclePause, defaultCyclePollPause)
}

func (c Channel) GetRegisterPause() time.Duration {
	return durationOrDefault(c.RegisterPause, defaultRegisterPollPause)
}

func (c Channel) findDeviceByTitle(title string) (*Device, error) {
	for _, v := range c.Devices {
		if v.Title == title {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("no device found for title '%s'", title)
}
