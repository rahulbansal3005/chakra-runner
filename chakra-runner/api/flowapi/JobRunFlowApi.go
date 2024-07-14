package api

import (
	"chakrarunner/api"
	"chakrarunner/client"
	"chakrarunner/pojo"
	"chakrarunner/util"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func RunFixedDelayJobs(DbConnection *gorm.DB, serverId string) {
	fixedDelayRunnableJobsQuery := "select j.* from jobs j inner join job_params jp on jp.job_id = j.id where j.server_id = \"" + serverId + "\" and jp.param_key = \"FIXED_DELAY\"  and TIMESTAMPDIFF(SECOND, j.end_time, NOW()) >= jp.param_value and j.is_enabled = true and j.status = \"STOPPED\" ;"
	runJobs(DbConnection, fixedDelayRunnableJobsQuery)
}

func RunFixedRateJobs(DbConnection *gorm.DB, serverId string) {
	fixedRateRunnableJobsQuery := "select j.* from jobs j inner join job_params jp on jp.job_id = j.id where j.server_id = \"" + serverId + "\" and jp.param_key = \"FIXED_RATE\"  and TIMESTAMPDIFF(SECOND, j.start_time, NOW()) >= jp.param_value and j.is_enabled = true;"
	runJobs(DbConnection, fixedRateRunnableJobsQuery)
}

func RunTimedOutJobs(DbConnection *gorm.DB, serverId string) {
	timeOutJobsQuery := "select j.* from jobs j inner join job_params jp on jp.job_id = j.id where j.server_id = \"" + serverId + "\" and jp.param_key = \"FIXED_DELAY\"  and TIMESTAMPDIFF(SECOND, j.start_time, NOW()) >= j.time_out and j.is_enabled = true and j.status = \"RUNNING\";"
	runJobs(DbConnection, timeOutJobsQuery)
}

func runJobs(DbConnection *gorm.DB, query string) {
	var jobs []pojo.Jobs
	DbConnection.Raw(query).Scan(&jobs)
	for i := 0; i < len(jobs); i++ {
		go runJob(DbConnection, &jobs[i])
	}
}

func runJob(DbConnection *gorm.DB, job *pojo.Jobs) {
	jobMetaData := api.GetJobMetaData(DbConnection, job.ID)
	orgCredentials := api.GetOrgCredentials(DbConnection, job.OrgId)

	err := jobStartedUpdate(DbConnection, job)

	if err != nil {
		util.Log.Error("Error starting job: " + job.JobName)
		return
	}
	util.Log.Info("STARTED Job: " + job.JobName + " " + strconv.Itoa(job.OrgId))

	job = api.GetJob(DbConnection, job.ID)

	err = client.TriggerJobRunCall(*job, *jobMetaData, *orgCredentials)

	if err == nil {
		jobCompleteUpdate(DbConnection, job)
	} else {
		jobFailedUpdate(DbConnection, job)
	}
	util.Log.Info("STOPPED Job: " + job.JobName + " " + strconv.Itoa(job.OrgId))

}

func jobStartedUpdate(DbConnection *gorm.DB, job *pojo.Jobs) error {
	return DbConnection.Transaction(func(tx *gorm.DB) error {
		job.StartTime = time.Now()
		job.Status = pojo.RUNNING
		err := api.UpdateJob(tx, job)
		if err != nil {
			return err
		}

		err = api.SaveAudit(tx, job.ID, job.ServerId, time.Now(), pojo.STARTED)
		if err != nil {
			return err
		}
		return nil
	})
}

func jobCompleteUpdate(DbConnection *gorm.DB, job *pojo.Jobs) {
	DbConnection.Transaction(func(tx *gorm.DB) error {

		job.EndTime = time.Now()
		job.Status = pojo.STOPPED
		err := api.UpdateJob(tx, job)

		err = api.SaveAudit(tx, job.ID, job.ServerId, time.Now(), pojo.SUCCESS)
		if err != nil {
			return err
		}
		return nil
	})
}

func jobFailedUpdate(DbConnection *gorm.DB, job *pojo.Jobs) {
	DbConnection.Transaction(func(tx *gorm.DB) error {

		job.Status = pojo.STOPPED
		job.EndTime = time.Now()
		err := api.UpdateJob(tx, job)
		if err != nil {
			return err
		}

		err = api.SaveAudit(tx, job.ID, job.ServerId, time.Now(), pojo.FAILURE)
		if err != nil {
			return err
		}
		return nil
	})
}
