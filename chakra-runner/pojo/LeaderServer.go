package pojo

import (
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type LeaderServer struct {
	gorm.Model
	ID       int
	ServerId string
	Version  optimisticlock.Version
}
