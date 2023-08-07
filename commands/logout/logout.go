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

	logoutRes, err := apiClient.ApiCall("logout", make(map[string]interface{}), apiClient.GetSessionID(), true, apiClient.IsProxyUsed())
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
