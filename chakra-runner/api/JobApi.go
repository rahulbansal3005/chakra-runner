package api

import (
	"chakrarunner/model"
	"chakrarunner/pojo"
	"chakrarunner/util"
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func SaveJob(db *gorm.DB, jobData model.JobData) (int, error) {
	jobPojo := &pojo.Jobs{JobName: jobData.JobName, StartTime: time.Now(), EndTime: time.Now(),
		Status: pojo.STOPPED, Url: jobData.URL, OrgId: jobData.OrgId, IsEnabled: jobData.IsEnabled, TimeOut: jobData.TimeOut}
	err := db.Create(jobPojo).Error
	if err != nil {
		return 0, err
	}
	return jobPojo.ID, nil
}

func SaveJobParams(db *gorm.DB, jobId int, jobData model.JobData) error {
	var jobParam = &pojo.JobParams{}
	if jobData.DelayType == string(pojo.FIXED_DELAY) || jobData.DelayType == string(pojo.FIXED_RATE) {
		jobParam = &pojo.JobParams{JobId: jobId, ParamKey: jobData.DelayType, ParamValue: strconv.Itoa(jobData.Frequency)}
	} else {
		jobParam = &pojo.JobParams{JobId: jobId, ParamKey: jobData.DelayType, ParamValue: jobData.Time}
	}
	return db.Create(jobParam).Error
}

func SaveJobMetaData(db *gorm.DB, jobId int, metaDataMap map[string]string) error {
	jobMetaDataList := make([]pojo.JobMetaData, 0, len(metaDataMap))
	for k := range metaDataMap {
		jobMetaDataList = append(jobMetaDataList, pojo.JobMetaData{ParamKey: k, ParamValue: metaDataMap[k], JobId: jobId})
	}
	return db.CreateInBatches(jobMetaDataList, len(metaDataMap)).Error
}

func DeleteJobMetaData(db *gorm.DB, jobId int) error {
	return db.Where("job_id = ?", jobId).Delete(&pojo.JobMetaData{}).Error
}

func GetJobServerCountMap(db *gorm.DB) []model.JobCountMapStruct {
	query := "select server_id, count(server_id) as count from jobs where server_id != \"\" group by server_id order by count(server_id) ;"
	var jobServerCountList []model.JobCountMapStruct
	tx := db.Raw(query)
	tx.Scan(&jobServerCountList)
	return jobServerCountList
}

func GetJob(DbConnection *gorm.DB, jobId int) *pojo.Jobs {
	jobPojo := &pojo.Jobs{}
	DbConnection.Where("id = ?", jobId).Find(jobPojo)
	return jobPojo
}

func GetJobIdsForServer(DbConnection *gorm.DB, serverId string, limit int) *[]pojo.Jobs {
	jobs := &[]pojo.Jobs{}
	DbConnection.Where("server_id = ?", serverId).Limit(limit).Find(jobs)
	return jobs
}

func GetJobParams(DbConnection *gorm.DB, jobId int) *pojo.JobParams {
	jobParamPojos := &pojo.JobParams{}
	DbConnection.Where("job_id = ?", jobId).Find(jobParamPojos)
	return jobParamPojos
}

func GetJobUsingNameAndOrg(db *gorm.DB, jobName string, orgId int) *pojo.Jobs {
	job := &pojo.Jobs{}
	db.Where("job_name = ? AND org_id >= ?", jobName, orgId).First(job)
	return job
}

func GetJobMetaData(DbConnection *gorm.DB, jobId int) *[]pojo.JobMetaData {
	jmd := &[]pojo.JobMetaData{}
	DbConnection.Where("job_id = ?", jobId).Find(jmd)
	return jmd
}

func GetTotalAssignedJobs(db *gorm.DB) int {
	getAssignedJobsQuery := "select count(*) from jobs where status != \"STOPPED\""
	var jobCount int
	db.Raw(getAssignedJobsQuery).Scan(&jobCount)
	return jobCount
}

func UpdateJob(DbConnection *gorm.DB, jobPojo *pojo.Jobs) error {
	return checkUpdateErrors(DbConnection.Updates(jobPojo), jobPojo.JobName, "job")
}

func UpdateJobParams(DbConnection *gorm.DB, jobParamPojo *pojo.JobParams) error {
	return checkUpdateErrors(DbConnection.Updates(jobParamPojo), strconv.Itoa(jobParamPojo.JobId), "jobParam")
}

func checkUpdateErrors(tx *gorm.DB, jobName string, pojoName string) error {
	if tx.RowsAffected == 0 {
		util.Log.Error("Error updating " + pojoName + jobName + ", row affected: " + strconv.Itoa(int(tx.RowsAffected)))
		return errors.New("Error updating " + pojoName + jobName)
	}
	if tx.Error != nil {
		util.Log.Error("Error updating " + pojoName + jobName + " Error:" + tx.Error.Error())
		return errors.New("Error updating" + pojoName + jobName)
	}
	return nil
}
