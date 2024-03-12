package model

import (
	"github.com/xhit/go-str2duration/v2"
	"time"
)

func durationOrDefault(input *string, defaultDuration time.Duration) time.Duration {
	if input == nil {
		return defaultDuration
	} else {
		v, err := str2duration.ParseDuration(*input)
		if err != nil {
			return defaultDuration
		}
		return v
	}
}
