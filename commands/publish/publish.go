package main

import (
	"fmt"
	"os"

	api_go_sdk "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	commands "github.com/mkol5222/my-cp-mgmt-commands/commands"
)

func takeOverSession(client api_go_sdk.ApiClient, session string) error {

	payload := map[string]interface{}{
		"uid": session,
	}

	takeOverSessionRes, err := client.ApiCall("take-over-session", payload, client.GetSessionID(), true, client.IsProxyUsed())
	if err != nil || !takeOverSessionRes.Success {

		if takeOverSessionRes.ErrorMsg != "" {
			fmt.Printf("take-over-session error: %s\n", fmt.Errorf(takeOverSessionRes.ErrorMsg))
			return fmt.Errorf(takeOverSessionRes.ErrorMsg)
		}
		fmt.Printf("take-over-session error\n")
		return fmt.Errorf(err.Error())
	} else {
		fmt.Println("take-over-session finished successfully")
	}
	return nil
}

func main() {
	apiClient, err := commands.InitClient()
	if err != nil {
		fmt.Println("Publish error: " + err.Error())
		os.Exit(1)
	}

	payload := make(map[string]interface{})
	if len(os.Args) < 2 {
		payload["uid"] = apiClient.GetSessionID()

	} else {
		payload["uid"] = os.Args[1]
	}

	takeOverSession(apiClient, fmt.Sprintf("%v", payload["uid"]))

	fmt.Printf("Publishing session %s\n", payload["uid"])

	publishRes, err := apiClient.ApiCall("publish", payload, apiClient.GetSessionID(), true, apiClient.IsProxyUsed())
	if err != nil {
		fmt.Println("Publish error: " + err.Error())
		os.Exit(1)
	}

	taskId := commands.ResolveTaskId(publishRes.GetData())

	if !publishRes.Success {
		errMsg := fmt.Sprintf("Publish failed: %s.", publishRes.ErrorMsg)
		if taskId != nil {
			errMsg += fmt.Sprintf(" task-id [%s]", taskId)
		}
		fmt.Println(errMsg)
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("Publish finished successfully. task-id [%s]", taskId))
}
