package api

import (
	"chakrarunner/api"
	"chakrarunner/config"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func LeaderSelectionCron() {
	fmt.Println(config.ServerId)
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Seconds().Do(func() {
		SendHeartBeat(config.DbConnection, config.ServerId)
		CheckLeaderHeartBeat(config.DbConnection, config.ServerId, config.ServerActiveThreshold)
	})

	s.StartBlocking()
}

func JobAssignmentCron() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Seconds().Do(func() {
		leaderServer := api.GetLeaderServer(config.DbConnection)
		if leaderServer.ServerId != config.ServerId {
			return
		}
		AssignJobs(config.DbConnection, config.ServerActiveThreshold)
	})
	s.StartBlocking()
}

func ActiveAndInActiveServerJobRedistribution() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(10).Seconds().Do(func() {
		leaderServer := api.GetLeaderServer(config.DbConnection)
		if leaderServer.ServerId != config.ServerId {
			return
		}
		unAssignUnActiveServerJobs(config.DbConnection, config.ServerActiveThreshold)
		AssignJobsToNewServer(config.DbConnection, config.ServerActiveThreshold)
	})
	s.StartBlocking()
}

func FixedDelayJobExecutionCron() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Seconds().Do(func() {
		RunFixedDelayJobs(config.DbConnection, config.ServerId)
	})

	s.StartBlocking()
}

func FixedRateJobExecutionCron() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Seconds().Do(func() {
		RunFixedRateJobs(config.DbConnection, config.ServerId)
	})

	s.StartBlocking()
}

func TimedOutJobExecutionCron() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Seconds().Do(func() {
		RunTimedOutJobs(config.DbConnection, config.ServerId)
	})

	s.StartBlocking()
}

func ChakraEventSyncCron() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(3).Seconds().Do(func() {
		leaderServer := api.GetLeaderServer(config.DbConnection)
		if leaderServer.ServerId == config.ServerId {
			syncEventsFromChakra(config.DbConnection)
		}
	})

	s.StartBlocking()
}
