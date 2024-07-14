package pojo

import (
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type ChakraSettings struct {
	gorm.Model
	ID                    int
	MinId                 int
	JobAssignmentRequired bool
	Version               optimisticlock.Version
}
