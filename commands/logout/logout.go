package main

import (
	"fmt"
	"os"

	commands "github.com/mkol5222/my-cp-mgmt-commands/commands"
)

func main() {
	apiClient, err := commands.InitClient()
	if err != nil {
		fmt.Println("logout error: " + err.Error())
		os.Exit(1)
	}

	payload := make(map[string]interface{})
	if len(os.Args) < 2 {
		payload["uid"] = apiClient.GetSessionID()
	} else {
		payload["uid"] = os.Args[1]
	}

	fmt.Printf("Logging out session %s\n", payload["uid"])

	logoutRes, err := apiClient.ApiCall("logout", payload, apiClient.GetSessionID(), true, apiClient.IsProxyUsed())
	if err != nil {
		fmt.Println("logout error: " + err.Error())
		os.Exit(1)
	}

	if !logoutRes.Success {
		errMsg := fmt.Sprintf("logOut failed: %s.", logoutRes.ErrorMsg)
		fmt.Println(errMsg)
		os.Exit(1)
	}
	fmt.Println("logout finished successfully")
}
