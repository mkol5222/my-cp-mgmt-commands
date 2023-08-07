package main

import (
	"fmt"
	"log"
	"os"

	api_go_sdk "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	commands "github.com/mkol5222/my-cp-mgmt-commands/commands"

	"github.com/TylerBrock/colorjson"
)

func showPackages(client api_go_sdk.ApiClient) error {

	payload := map[string]interface{}{
		"details-level": "full",
	}

	showPackagesRes, err := client.ApiCall("show-packages", payload, client.GetSessionID(), true, client.IsProxyUsed())
	if err != nil || !showPackagesRes.Success {
		if showPackagesRes.ErrorMsg != "" {
			return fmt.Errorf(showPackagesRes.ErrorMsg)
		}
		return fmt.Errorf(err.Error())
	}

	//log.Printf("show-Packages response: %v\n", showPackagesRes)

	showPackagesJson := showPackagesRes.GetData()

	objectsList, ok := showPackagesJson["packages"].([]interface{})

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
			f := colorjson.NewFormatter()
			f.Indent = 4

			// Marshall the Colorized JSON
			s, _ := f.Marshal(objectsMap)
			fmt.Println(string(s))

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

	showPackages(apiClient)
}
