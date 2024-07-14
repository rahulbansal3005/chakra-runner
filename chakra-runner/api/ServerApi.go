package api

import (
	"chakrarunner/pojo"
	"time"

	"gorm.io/gorm"
)

func SaveHeartBeat(db *gorm.DB, serverId string) {
	serverPojo := &pojo.ServerHeartBeat{ServerId: serverId, LastHeartBeat: time.Now()}
	db.Create(serverPojo)
}

func UpdateHeartBeat(db *gorm.DB, serverId string) {
	heartBeat := GetServerHeartBeat(db, serverId)
	heartBeat.LastHeartBeat = time.Now()
	db.Updates(heartBeat)
}

func GetServerHeartBeat(db *gorm.DB, serverId string) *pojo.ServerHeartBeat {
	hb := &pojo.ServerHeartBeat{}
	db.Where("server_id = ?", serverId).Find(&hb)
	return hb
}

func GetLeaderServer(db *gorm.DB) *pojo.LeaderServer {
	leaderServer := &pojo.LeaderServer{}
	db.First(&leaderServer)
	return leaderServer
}

func GetActiveServers(db *gorm.DB, serverActivethreshold int) []string {
	var servers = &[]pojo.ServerHeartBeat{}
	db.Find(&servers)

	var activeServers []string
	for _, s := range *servers {
		if float64(time.Since(s.LastHeartBeat).Seconds()) < float64(serverActivethreshold) {
			activeServers = append(activeServers, s.ServerId)
		}
	}
	return activeServers
}

func GetInActiveServers(db *gorm.DB, serverActivethreshold int) []string {
	var servers = &[]pojo.ServerHeartBeat{}
	db.Find(&servers)

	var activeServers []string
	for _, s := range *servers {
		if float64(time.Since(s.LastHeartBeat).Seconds()) > float64(serverActivethreshold) {
			activeServers = append(activeServers, s.ServerId)
		}
	}
	return activeServers
}
