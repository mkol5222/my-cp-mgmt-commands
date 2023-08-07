package main

import (
	"fmt"
	"log"
	"os"

	api_go_sdk "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	commands "github.com/mkol5222/my-cp-mgmt-commands/commands"
)

func showSessions(client api_go_sdk.ApiClient) error {

	payload := map[string]interface{}{
		"details-level": "full",
	}

	showSessionsRes, err := client.ApiCall("show-sessions", payload, client.GetSessionID(), true, client.IsProxyUsed())
	if err != nil || !showSessionsRes.Success {
		if showSessionsRes.ErrorMsg != "" {
			return fmt.Errorf(showSessionsRes.ErrorMsg)
		}
		return fmt.Errorf(err.Error())
	}

	//log.Printf("show-sessions response: %v\n", showSessionsRes)

	showSessionsJson := showSessionsRes.GetData()

	//log.Printf("show-sessions json: %v\n", showSessionsJson)

	objectsList, ok := showSessionsJson["objects"].([]interface{})

	if ok {
		for i := range objectsList {
			objectsMap := objectsList[i].(map[string]interface{})
			tempObject := make(map[string]interface{})

			if v, _ := objectsMap["name"]; v != nil {
				tempObject["name"] = v
			} else {
				tempObject["name"] = ""
			}

			if v, _ := objectsMap["changes"]; v != nil {
				tempObject["changes"] = v
			} else {
				tempObject["changes"] = 0
			}
			if v, _ := objectsMap["locks"]; v != nil {
				tempObject["locks"] = v
			} else {
				tempObject["locks"] = 0
			}

			if v, _ := objectsMap["uid"]; v != nil {
				tempObject["uid"] = v
			}

			fmt.Printf("- %s: L: %#v CH: %#v %s\n\n", tempObject["uid"], tempObject["locks"], tempObject["changes"], tempObject["name"])
			//fmt.Printf("%v\n", objectsMap)

		}
	} else {
		log.Printf("objectsList is not a list: %v\n", objectsList)
	}

	return nil
}

func main() {
	apiClient, err := commands.InitClient()
	if err != nil {
		fmt.Println("logout error: " + err.Error())
		os.Exit(1)
	}
	log.Printf("Api client initialized successfully with SID %v\n", apiClient.GetSessionID())

	showSessions(apiClient)
}
