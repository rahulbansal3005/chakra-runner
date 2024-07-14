package api

import (
	"chakrarunner/pojo"

	"gorm.io/gorm"
)

func SaveChakraSetting(db *gorm.DB, minId int) *pojo.ChakraSettings {
	chakraSetting := &pojo.ChakraSettings{MinId: minId, JobAssignmentRequired: false}
	db.Create(chakraSetting)
	return chakraSetting
}

func UpdateChakraSetting(DbConnection *gorm.DB, chakraSetting *pojo.ChakraSettings) error {
	return DbConnection.Updates(chakraSetting).Error
}

func GetChakraSetting(DbConnection *gorm.DB) *pojo.ChakraSettings {
	chakraSetting := &pojo.ChakraSettings{}
	DbConnection.Find(&chakraSetting)
	return chakraSetting
}
