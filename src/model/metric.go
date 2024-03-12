package model

import (
	"fmt"
	"time"
)

type Metric struct {
	Key       string    `json:"key"`
	Channel   string    `json:"channel"`
	Device    string    `json:"device"`
	Alias     string    `json:"alias"`
	Register  string    `json:"register"`
	RawValue  any       `json:"raw"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func (m Metric) IsExpired(ttl time.Duration) bool {
	return m.Timestamp.Add(ttl).Before(time.Now())
}
func (m Metric) String() string {
	return fmt.Sprintf("key: %-30s raw: %-10d val: %-10.2f ts: %s", m.Key, m.RawValue, m.Value, m.Timestamp.Format("2006-01-02 15:04:05.000"))
}
func MetricKey(register *Register) string {
	return fmt.Sprintf("%s:%s:%s", register.Device.Channel.Title, register.Device.Title, register.Title)
}
