package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DbConnection *gorm.DB

func ConnectDB() {
	var dsn = DbUsername + ":" + DbPassword + "@tcp(" + DbUrl
	DbConnection, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
