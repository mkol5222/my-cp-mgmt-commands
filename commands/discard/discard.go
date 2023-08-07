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

	discardRes, err := apiClient.ApiCall("discard", map[string]interface{}{}, apiClient.GetSessionID(), true, apiClient.IsProxyUsed())
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
