package pojo

import (
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type ParamKey string

const (
	FIXED_RATE  ParamKey = "FIXED_RATE"
	FIXED_DELAY ParamKey = "FIXED_DELAY"
	FIXED_TIME  ParamKey = "FIXED_TIME"
)

type JobParams struct {
	gorm.Model
	ID         int
	JobId      int
	ParamKey   string
	ParamValue string
	Version    optimisticlock.Version
}
