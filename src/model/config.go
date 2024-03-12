package model

import (
	"fmt"
	"strings"
	"time"
)

const defaultMetricTTL = time.Second * 30

type Config struct {
	Ttl              *string   `json:"ttl,omitempty"`
	PrometheusExport bool      `json:"export_prometheus,omitempty"`
	Channels         []Channel `json:"channels,omitempty"`
}

func (config *Config) GetTTL() time.Duration {
	return durationOrDefault(config.Ttl, defaultMetricTTL)
}
func (config *Config) FindChannelByTitle(title string) (*Channel, error) {
	for _, v := range config.Channels {
		if v.Title == title {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("no channel found for title '%s'", title)
}
func (config *Config) FindRegister(reference string) (*Register, error) {
	ref := strings.Split(reference, ":")
	if 3 != len(ref) {
		return nil, fmt.Errorf("invalid reference passed: '%s'", reference)
	}
	c, err := config.FindChannelByTitle(strings.TrimSpace(ref[0]))
	if nil != err {
		return nil, err
	}
	d, err := c.findDeviceByTitle(strings.TrimSpace(ref[1]))
	if nil != err {
		return nil, err
	}
	r, err := d.findRegisterByTitle(strings.TrimSpace(ref[2]))
	return r, err
}
