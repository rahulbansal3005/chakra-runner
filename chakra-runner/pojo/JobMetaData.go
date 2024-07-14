package pojo

import (
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type JobMetaData struct {
	gorm.Model
	ID         int
	JobId      int
	ParamKey   string
	ParamValue string
	Version    optimisticlock.Version
}
