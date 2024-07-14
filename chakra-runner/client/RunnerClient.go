package client

import (
	"bytes"
	"chakrarunner/pojo"
	"chakrarunner/util"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func TriggerJobRunCall(job pojo.Jobs, jobMetaData []pojo.JobMetaData, credentials []pojo.OrgCredentials) error {
	client := &http.Client{}

	req, err := http.NewRequest("POST", job.Url, bytes.NewBuffer(getRequestBody(jobMetaData, job.JobName)))
	if checkForError(err, "Generating Request for JobId: "+strconv.Itoa(job.ID)) {
		return err
	}

	setHeaders(req, credentials)
	resp, reqError := client.Do(req)
	if checkForError(reqError, "Error in Job, jobName: "+job.JobName+", orgId: "+strconv.Itoa(job.OrgId)) {
		return reqError
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		util.Log.Error("Error in Job, jobName: " + job.JobName + ", orgId: " + strconv.Itoa(job.OrgId) + " " + string(body))
		return errors.New("error while triggering")
	}

	return nil
}

func setHeaders(req *http.Request, credentials []pojo.OrgCredentials) {
	req.Header.Add("Content-Type", "application/json")
	for _, c := range credentials {
		req.Header.Add(c.HeaderParam, c.HeaderValue)
	}
}

func getRequestBody(jobMetaData []pojo.JobMetaData, jobName string) []byte {
	metaDataMap := make(map[string]string)

	for _, data := range jobMetaData {
		metaDataMap[data.ParamKey] = data.ParamValue
	}
	metaDataMap["JobName"] = jobName
	reqBody, err := json.Marshal(metaDataMap)

	if checkForError(err, "getRequestBody") {
		return nil
	}
	return reqBody
}

func checkForError(err error, message string) bool {
	if err != nil {
		util.Log.Error(message + " ,error: " + err.Error())
		return true
	}
	return false
}
