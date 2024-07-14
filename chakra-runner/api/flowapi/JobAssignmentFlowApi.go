package api

import (
	"chakrarunner/api"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func AssignJobs(DbConnection *gorm.DB, serverActiveThreshold int) {
	unassignedJobQuery := "select id from jobs where server_id = \"\" and is_enabled = true;"
	var jobIds []int
	DbConnection.Raw(unassignedJobQuery).Scan(&jobIds)
	if jobIds == nil {
		return
	}
	assignJobs(DbConnection, serverActiveThreshold, jobIds)
}

func unAssignUnActiveServerJobs(DbConnection *gorm.DB, serverActiveThreshold int) {
	inActiveServers := api.GetInActiveServers(DbConnection, serverActiveThreshold)
	unAssignInActiveServerQuery := "update jobs set server_id = \"\" where server_id in (\"" + strings.Join(inActiveServers, "\",\"") + "\");"
	DbConnection.Exec(unAssignInActiveServerQuery)
}

func AssignJobsToNewServer(DbConnection *gorm.DB, serverActiveThreshold int) {
	activeServers := api.GetActiveServers(DbConnection, serverActiveThreshold)
	jobServerCountList := api.GetJobServerCountMap(DbConnection)
	jobServerCountMap := make(map[string]int)
	var totalAssignedJobs = 0
	for i := range jobServerCountList {
		totalAssignedJobs += jobServerCountList[i].Count
		jobServerCountMap[jobServerCountList[i].ServerId] = jobServerCountList[i].Count
	}

	var freeServer string
	for _, server := range activeServers {
		v, present := jobServerCountMap[server]
		if !present || v == 0 {
			freeServer = server
			break
		}
	}

	if freeServer == "" {
		return
	}

	idealJobsPerServer := totalAssignedJobs / len(activeServers)
	for k, _ := range jobServerCountMap {
		limit := jobServerCountMap[k] - idealJobsPerServer
		jobs := api.GetJobIdsForServer(DbConnection, k, limit)
		var jobIds []string
		for _, j := range *jobs {
			jobIds = append(jobIds, strconv.Itoa(j.ID))
		}
		updateQuery := "update jobs set server_id = \"" + freeServer + "\" where id in (" + strings.Join(jobIds, ",") + ");"
		DbConnection.Exec(updateQuery)
	}

}

func assignJobs(DbConnection *gorm.DB, serverActiveThreshold int, jobIds []int) {
	activeServers := api.GetActiveServers(DbConnection, serverActiveThreshold)
	jobServerCountList := api.GetJobServerCountMap(DbConnection)

	jobServerCountMap := make(map[string]int)

	for i := range jobServerCountList {
		jobServerCountMap[jobServerCountList[i].ServerId] = jobServerCountList[i].Count
	}

	var maxJobsAssigned = 0

	for _, v := range jobServerCountMap {
		if v > maxJobsAssigned {
			maxJobsAssigned = v
		}
	}

	serverToJobAssignmentMap := make(map[string][]string)

	var totalAssignedJobs int = 0
	var totalJobsToAssign int = len(jobIds)
	for _, serverId := range activeServers {
		var serverJobLimit int = maxJobsAssigned - jobServerCountMap[serverId]
		for serverJobLimit > 0 && totalAssignedJobs < totalJobsToAssign {
			serverToJobAssignmentMap[serverId] = append(serverToJobAssignmentMap[serverId], strconv.Itoa(jobIds[totalAssignedJobs]))
			serverJobLimit -= 1
			totalAssignedJobs += 1
		}
	}

	if totalJobsToAssign != totalAssignedJobs {
		for _, serverId := range activeServers {
			var serverJobLimit int = 5
			for serverJobLimit > 0 && totalAssignedJobs < totalJobsToAssign {
				serverToJobAssignmentMap[serverId] = append(serverToJobAssignmentMap[serverId], strconv.Itoa(jobIds[totalAssignedJobs]))
				serverJobLimit -= 1
				totalAssignedJobs += 1
			}
		}
	}

	for k := range serverToJobAssignmentMap {
		updateQuery := "update jobs set server_id = \"" + k + "\" where id in (" + strings.Join(serverToJobAssignmentMap[k], ",") + ");"
		DbConnection.Exec(updateQuery)
	}
}
