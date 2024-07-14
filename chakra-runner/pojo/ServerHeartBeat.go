package pojo

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type ServerHeartBeat struct {
	gorm.Model
	ID            int
	ServerId      string
	LastHeartBeat time.Time
	Version       optimisticlock.Version
}
