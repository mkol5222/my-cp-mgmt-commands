package main

import (
	"fmt"
	"os"

	commands "github.com/mkol5222/my-cp-mgmt-commands/commands"
)

func main() {
	apiClient, err := commands.InitClient()
	if err != nil {
		fmt.Println("Discard error: " + err.Error())
		os.Exit(1)
	}

	payload := make(map[string]interface{})
	if len(os.Args) < 2 {
		payload["uid"] = apiClient.GetSessionID()
	} else {
		payload["uid"] = os.Args[1]
	}

	fmt.Printf("Discarding session %s\n", payload["uid"])

	discardRes, err := apiClient.ApiCall("discard", payload, apiClient.GetSessionID(), true, apiClient.IsProxyUsed())
	if err != nil {
		fmt.Println("Discard error: " + err.Error())
		os.Exit(1)
	}

	if !discardRes.Success {
		errMsg := fmt.Sprintf("Discard failed: %s.", discardRes.ErrorMsg)
		fmt.Println(errMsg)
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("Discard finished successfully."))
}
