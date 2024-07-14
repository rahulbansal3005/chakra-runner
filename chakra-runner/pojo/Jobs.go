package pojo

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type Status string

const (
	STOPPED Status = "STOPPED"
	RUNNING Status = "RUNNING"
)

type Jobs struct {
	gorm.Model
	ID        int
	JobName   string
	StartTime time.Time
	EndTime   time.Time
	Status    Status
	Url       string
	OrgId     int
	ServerId  string
	IsEnabled bool
	TimeOut   int
	Version   optimisticlock.Version
}
