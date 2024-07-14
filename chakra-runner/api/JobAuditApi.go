package api

import (
	"chakrarunner/pojo"
	"chakrarunner/util"
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func SaveAudit(db *gorm.DB, jobId int, serverId string, timeStamp time.Time, resolution pojo.Resolution) error {
	auditPojo := &pojo.JobAudit{ServerId: serverId, Time: timeStamp, JobId: jobId, Resolution: resolution}
	err := db.Create(auditPojo).Error
	if err != nil {
		util.Log.Error(err)
		return errors.New("Error Saving Audit for jobId: " + strconv.Itoa(jobId) + " " + err.Error())
	}
	return nil
}
