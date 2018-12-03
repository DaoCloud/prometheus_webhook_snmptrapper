package types

import (
	"time"
)

type AlertGroup struct {
	Status    string
	Alerts    []Alert
}

type Alert struct {
	Status       string
	Annotations  map[string]string
	Labels       map[string]string
	StartsAt     time.Time
	EndsAt       time.Time
	GeneratorURL string
}
