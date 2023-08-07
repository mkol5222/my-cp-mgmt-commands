package main

import (
	"fmt"
	"os"

	commands "github.com/mkol5222/my-cp-mgmt-commands/commands"
)

func main() {
	apiClient, err := commands.InitClient()
	if err != nil {
		fmt.Println("Approve Session error: " + err.Error())
		os.Exit(1)
	}
	if len(os.Args) < 2 {
		fmt.Println("no session uid provided")
		os.Exit(1)
	}
	payload := make(map[string]interface{})
	payload["uid"] = os.Args[1]
	approveSessionRes, err := apiClient.ApiCall("approve-session", payload, apiClient.GetSessionID(), true, apiClient.IsProxyUsed())
	if err != nil {
		fmt.Println("Approve Session error: " + err.Error())
		os.Exit(1)
	}

	if !approveSessionRes.Success {
		errMsg := fmt.Sprintf("Approve Session failed: %s.", approveSessionRes.ErrorMsg)
		fmt.Println(errMsg)
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("Approve Session finished successfully."))
}
