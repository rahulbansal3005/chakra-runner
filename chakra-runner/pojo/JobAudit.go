package pojo

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type Resolution string

const (
	SUCCESS Resolution = "SUCCESS"
	FAILURE            = "FAILURE"
	TIMEOUT            = "TIMEOUT"
	STARTED            = "STARTED"
)

type JobAudit struct {
	gorm.Model
	ID         int
	JobId      int
	Time       time.Time
	ServerId   string
	Resolution Resolution
	Version    optimisticlock.Version
}
