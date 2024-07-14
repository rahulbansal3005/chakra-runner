package client

import (
	"chakrarunner/config"
	"chakrarunner/util"
	"io"
	"net/http"
	"strconv"
)

func SyncEvents(minId int) []byte {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", config.ChakraBaseUrl+"/job-ledger", nil)

	setChakraHeaders(req)
	q := req.URL.Query()
	q.Add("minId", strconv.Itoa(minId))
	req.URL.RawQuery = q.Encode()
	resp, reqError := client.Do(req)
	if checkForError(reqError, "Error syncing events from chakra, minId: "+strconv.Itoa(minId)) {
		return nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	util.Log.Info("Sync event done for minId: " + strconv.Itoa(minId))
	return body
}

func setChakraHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authDomainName", config.ChakraAuthDomainName)
	req.Header.Add("authUsername", config.ChakraAuthUserName)
	req.Header.Add("authPassword", config.ChakraAuthPassword)

}
