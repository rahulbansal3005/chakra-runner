package api

import (
	"chakrarunner/api"
	"chakrarunner/client"
	"chakrarunner/model"
	"chakrarunner/pojo"
	"chakrarunner/util"
	"encoding/json"
	"strconv"

	"gorm.io/gorm"
)

func syncEventsFromChakra(DbConnection *gorm.DB) {
	chakraSetting := api.GetChakraSetting(DbConnection)

	if chakraSetting.ID == 0 {
		chakraSetting = api.SaveChakraSetting(DbConnection, 0)
	}
	var eventLedgerData model.JobLedgerSyncData
	json.Unmarshal(client.SyncEvents(chakraSetting.MinId), &eventLedgerData)

	for _, s := range eventLedgerData.Events {
		var err error
		if s.EventType == "CREATE" {
			job := api.GetJobUsingNameAndOrg(DbConnection, s.JobData.JobName, s.JobData.OrgId)
			if job.ID == 0 {
				err = saveJobData(DbConnection, s.JobData)
			}
			chakraSetting.JobAssignmentRequired = true
		}

		if s.EventType == "UPDATE" {
			job := api.GetJobUsingNameAndOrg(DbConnection, s.JobData.JobName, s.JobData.OrgId)
			if job.ID == 0 {
				util.Log.Error("Job: " + s.JobData.JobName + ", " + strconv.Itoa(s.JobData.OrgId) + " not present so cannot update")
			} else {
				err = updateJobData(DbConnection, s.JobData, job, chakraSetting)
			}
		}

		if err != nil {
			break
		}

		err = saveOrgCredentials(DbConnection, s.JobData.OrgCredentials, s.JobData.OrgId)
		if err != nil {
			util.Log.Error(err.Error())
			return
		}
		if s.ID > chakraSetting.MinId {
			chakraSetting.MinId = s.ID
		}
	}

	api.UpdateChakraSetting(DbConnection, chakraSetting)
}

func saveJobData(DbConnection *gorm.DB, data model.JobData) error {
	return DbConnection.Transaction(func(tx *gorm.DB) error {

		jobId, err := api.SaveJob(tx, data)
		if err != nil {
			return err
		}

		err = api.SaveJobParams(tx, jobId, data)
		if err != nil {
			return err
		}

		err = api.SaveJobMetaData(tx, jobId, data.JobMetaData)
		if err != nil {
			return err
		}
		return nil
	})
}

func updateJobData(DbConnection *gorm.DB, data model.JobData, job *pojo.Jobs, chakraSetting *pojo.ChakraSettings) error {
	return DbConnection.Transaction(func(tx *gorm.DB) error {
		if job.IsEnabled != data.IsEnabled {
			job.IsEnabled = data.IsEnabled
			chakraSetting.JobAssignmentRequired = true
			job.ServerId = ""
		}

		job.Url = data.URL
		job.TimeOut = data.TimeOut
		err := api.UpdateJob(tx, job)
		if err != nil {
			return err
		}

		if job.IsEnabled == false {
			err = tx.Exec("update jobs set is_enabled = false, server_id = \"\" where id = " + strconv.Itoa(job.ID) + ";").Error
			if err != nil {
				util.Log.Error(err.Error())
				return err
			}
		}

		updatedParams := api.GetJobParams(tx, job.ID)
		if data.DelayType == string(pojo.FIXED_DELAY) || data.DelayType == string(pojo.FIXED_RATE) {
			updatedParams.ParamValue = strconv.Itoa(data.Frequency)
		} else {
			updatedParams.ParamValue = data.Time
		}
		err = api.UpdateJobParams(tx, updatedParams)
		if err != nil {
			return err
		}

		err = api.DeleteJobMetaData(tx, job.ID)
		if err != nil {
			return err
		}

		err = api.SaveJobMetaData(tx, job.ID, data.JobMetaData)
		if err != nil {
			return err
		}

		return nil
	})
}

func saveOrgCredentials(DbConnection *gorm.DB, credentials map[string]string, orgId int) error {
	err := DbConnection.Transaction(func(tx *gorm.DB) error {

		err := api.DeleteOrgCredentials(tx, orgId)
		if err != nil {
			return err
		}

		err = api.SaveOrgCredentials(tx, orgId, credentials)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}
