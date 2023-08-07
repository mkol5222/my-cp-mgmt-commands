package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
)

const (
	DefaultFilename = "sid.json"
)

type Session struct {
	Sid string `json:"sid"`
	Uid string `json:"uid"`
}

func (s *Session) Save(sessionFileName string) error {
	f, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(sessionFileName, f, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetSession(sessionFileName string) (Session, error) {
	if _, err := os.Stat(sessionFileName); os.IsNotExist(err) {
		_, err := os.Create(sessionFileName)
		if err != nil {
			return Session{}, err
		}
	}
	b, err := ioutil.ReadFile(sessionFileName)
	if err != nil || len(b) == 0 {
		return Session{}, err
	}
	var s Session
	if err = json.Unmarshal(b, &s); err != nil {
		return Session{}, err
	}
	return s, nil
}

func ResolveTaskId(data map[string]interface{}) interface{} {
	if data != nil {
		if v := data["tasks"]; v != nil {
			tasks := v.([]interface{})
			if len(tasks) > 0 {
				return tasks[0].(map[string]interface{})["task-id"]
			}
		}

		if v := data["task-id"]; v != nil {
			return v
		}
	}
	return nil
}

func login(client *checkpoint.ApiClient, username string, pwd string, apiKey string, domain string, sessionName string, sessionDescription string, sessionTimeout int) (Session, error) {
	log.Printf("Perform login")
	var loginRes checkpoint.APIResponse
	var err error

	payload := make(map[string]interface{})
	if sessionName != "" {
		payload["session-name"] = sessionName
	}

	if sessionDescription != "" {
		payload["session-description"] = sessionDescription
	}

	if sessionTimeout != -1 {
		payload["session-timeout"] = sessionTimeout
	}

	if apiKey != "" {
		loginRes, err = client.ApiLoginWithApiKey(apiKey, false, domain, false, payload)
	} else {
		loginRes, err = client.ApiLogin(username, pwd, false, domain, false, payload)
	}

	if err != nil {
		localRequestsError := "invalid character '<' looking for beginning of value"
		if strings.Contains(err.Error(), localRequestsError) {
			return Session{}, fmt.Errorf("login failure: API server needs to be configured to accept requests from all IP addresses")
		}
		return Session{}, err
	}
	if !loginRes.Success {
		return Session{}, fmt.Errorf(loginRes.ErrorMsg)
	}

	uid := ""
	if val, ok := loginRes.GetData()["uid"]; ok {
		uid = val.(string)
	}

	s := Session{
		Sid: client.GetSessionID(),
		Uid: uid,
	}

	return s, nil
}

func InitClient() (checkpoint.ApiClient, error) {
	// Default values
	port := checkpoint.DefaultPort
	timeout := checkpoint.TimeOut
	proxyPort := checkpoint.DefaultProxyPort
	autoPublishBatchSize := -1

	// Get credentials from Environment variables
	server := os.Getenv("CHECKPOINT_SERVER")
	username := os.Getenv("CHECKPOINT_USERNAME")
	password := os.Getenv("CHECKPOINT_PASSWORD")
	portVal := os.Getenv("CHECKPOINT_PORT")
	timeoutVal := os.Getenv("CHECKPOINT_TIMEOUT")
	sessionFileName := os.Getenv("CHECKPOINT_SESSION_FILE_NAME")
	proxyHost := os.Getenv("CHECKPOINT_PROXY_HOST")
	proxyPortStr := os.Getenv("CHECKPOINT_PROXY_PORT")
	apiKey := os.Getenv("CHECKPOINT_API_KEY")
	cloudMgmtId := os.Getenv("CHECKPOINT_CLOUD_MGMT_ID")
	autoPublishBatchSizeVal := os.Getenv("CHECKPOINT_AUTO_PUBLISH_BATCH_SIZE")

	var err error
	if portVal != "" {
		port, err = strconv.Atoi(portVal)
		if err != nil {
			return checkpoint.ApiClient{}, fmt.Errorf("failed to parse CHECKPOINT_PORT to integer")
		}
	}

	if proxyPortStr != "" {
		proxyPort, err = strconv.Atoi(proxyPortStr)
		if err != nil {
			return checkpoint.ApiClient{}, fmt.Errorf("failed to parse CHECKPOINT_PROXY_PORT to integer")
		}
	}

	if timeoutVal != "" {
		timeoutInteger, err := strconv.Atoi(timeoutVal)
		if err != nil {
			return checkpoint.ApiClient{}, fmt.Errorf("failed to parse CHECKPOINT_TIMEOUT to integer")
		}
		timeout = time.Duration(timeoutInteger)
	}

	if sessionFileName == "" {
		sessionFileName = DefaultFilename
	}

	if autoPublishBatchSizeVal != "" {
		autoPublishBatchSize, err = strconv.Atoi(timeoutVal)
		if err != nil {
			return checkpoint.ApiClient{}, fmt.Errorf("failed to parse CHECKPOINT_AUTO_PUBLISH_BATCH_SIZE to integer")
		}
	}

	if server == "" || ((username == "" || password == "") && apiKey == "") {
		return checkpoint.ApiClient{}, fmt.Errorf("missing at least one required parameter to initialize API client (CHECKPOINT_SERVER, (CHECKPOINT_USERNAME and CHECKPOINT_PASSWORD) OR CHECKPOINT_API_KEY)")
	}

	// install policy/publish - only on management api
	if val, ok := os.LookupEnv("CHECKPOINT_CONTEXT"); ok {
		if val == "gaia_api" {
			return checkpoint.ApiClient{}, fmt.Errorf("post apply/destroy scripts are valid only on management api. Env var CHECKPOINT_CONTEXT is 'gaia_api'")
		}
	}

	args := checkpoint.ApiClientArgs{
		Port:                    port,
		Fingerprint:             "",
		Sid:                     "",
		Server:                  server,
		ProxyHost:               proxyHost,
		ProxyPort:               proxyPort,
		ApiVersion:              "",
		IgnoreServerCertificate: false,
		AcceptServerCertificate: false,
		DebugFile:               "deb.txt",
		Context:                 "web_api",
		Timeout:                 timeout,
		Sleep:                   checkpoint.SleepTime,
		UserAgent:               "Terraform",
		CloudMgmtId:             cloudMgmtId,
		AutoPublishBatchSize:    autoPublishBatchSize,
	}

	s, err := GetSession(sessionFileName)
	if err != nil {
		return checkpoint.ApiClient{}, err
	}
	if s.Sid != "" {
		args.Sid = s.Sid
	} else {
		log.Printf("Session not found. Performing login")
		mgmtClient := checkpoint.APIClient(args)
		s, err = login(mgmtClient, username, password, apiKey, "", "", "", -1)
		if err == nil {
			// log.Printf("Login result: %v", s)
			log.Printf("Saving session to: %v", DefaultFilename)
			s.Save(DefaultFilename)
			args.Sid = s.Sid
		} else {
			log.Printf("Failed to login: %v", err)
			return checkpoint.ApiClient{}, fmt.Errorf("session id not found. Verify %s file exists in working directory", sessionFileName)
		}

	}

	mgmt := checkpoint.APIClient(args)

	return *mgmt, nil
}
