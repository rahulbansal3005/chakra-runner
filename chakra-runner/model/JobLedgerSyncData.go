package model

type JobLedgerSyncData struct {
	Events []Event `json:"events"`
}

type Event struct {
	ID        int     `json:"id"`
	EventType string  `json:"eventType"`
	JobData   JobData `json:"jobData"`
}

type JobData struct {
	JobName        string            `json:"jobName"`
	OrgId          int               `json:"orgId"`
	URL            string            `json:"url"`
	IsEnabled      bool              `json:"isEnabled"`
	DelayType      string            `json:"delayType"`
	Frequency      int               `json:"frequency"`
	Time           string            `json:"time"`
	TimeOut        int               `json:"timeOut"`
	JobMetaData    map[string]string `json:"jobMetaData"`
	OrgCredentials map[string]string `json:"orgCredentials"`
}
