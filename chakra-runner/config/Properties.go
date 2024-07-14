package config

import (
	"github.com/magiconair/properties"
)

var ServerId string
var DbUrl string
var DbUsername string
var DbPassword string
var ServerActiveThreshold int
var JobTimeOut int
var ChakraBaseUrl string
var ChakraAuthDomainName string
var ChakraAuthUserName string
var ChakraAuthPassword string

func PopulateProperties() {
	p := properties.MustLoadFile("chakra-runner.properties", properties.UTF8)
	ServerId = p.MustGetString("serverId")
	DbUrl = p.MustGetString("dbUrl")
	DbUsername = p.MustGetString("dbUsername")
	DbPassword = p.MustGetString("dbPassword")
	ServerActiveThreshold = p.GetInt("serverActiveThreshold", 5)
	JobTimeOut = p.GetInt("jobTimeOut", 420)
	ChakraBaseUrl = p.MustGetString("chakraBaseUrl")
	ChakraAuthDomainName = p.MustGetString("chakraAuthDomainName")
	ChakraAuthUserName = p.MustGetString("chakraAuthUserName")
	ChakraAuthPassword = p.MustGetString("chakraAuthPassword")
}
