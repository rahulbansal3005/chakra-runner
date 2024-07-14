package main

import (
	flowApi "chakrarunner/api/flowApi"
	"chakrarunner/config"
	"chakrarunner/util"
)

func main() {
	config.PopulateProperties()
	config.ConnectDB()
	// config.DbConnection.AutoMigrate(&pojo.JobAudit{}, &pojo.JobMetaData{}, &pojo.JobParams{}, &pojo.Jobs{}, &pojo.OrgCredentials{},
	// 	&pojo.ServerHeartBeat{}, &pojo.LeaderServer{}, &pojo.ChakraJobSyncDetail{})
	util.SetupLogger()
	go flowApi.LeaderSelectionCron()
	go flowApi.JobAssignmentCron()
	go flowApi.ActiveAndInActiveServerJobRedistribution()
	go flowApi.FixedRateJobExecutionCron()
	go flowApi.FixedDelayJobExecutionCron()
	go flowApi.TimedOutJobExecutionCron()
	go flowApi.ChakraEventSyncCron()
	select {}
}
