package pojo

import (
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type OrgCredentials struct {
	gorm.Model
	ID          int
	OrgId       int
	HeaderParam string
	HeaderValue string
	Version     optimisticlock.Version
}
