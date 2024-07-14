package api

import (
	"chakrarunner/api"
	"chakrarunner/util"
	"time"

	"gorm.io/gorm"
)

func SendHeartBeat(dbConnection *gorm.DB, serverId string) {
	heartBeat := api.GetServerHeartBeat(dbConnection, serverId)
	if heartBeat.ServerId == "" {
		api.SaveHeartBeat(dbConnection, serverId)
	} else {
		api.UpdateHeartBeat(dbConnection, serverId)
	}
}

func CheckLeaderHeartBeat(dbConnection *gorm.DB, serverId string, serverActiveThreashold int) {
	leaderServer := api.GetLeaderServer(dbConnection)
	leaderHB := api.GetServerHeartBeat(dbConnection, leaderServer.ServerId)
	if float64(time.Since(leaderHB.LastHeartBeat).Seconds()) > float64(serverActiveThreashold) {
		leaderServer.ServerId = serverId
		tx := dbConnection.Updates(leaderServer)
		if tx.Error != nil {
			util.Log.Error("Error: " + tx.Error.Error())
		}
		util.Log.Error("New Leader: " + serverId)
	}
}
