package main

import (
	"fmt"
	"log"
	"os"

	commands "github.com/mkol5222/my-cp-mgmt-commands/commands"
)

func main() {
	apiClient, err := commands.InitClient()
	if err != nil {
		fmt.Println("logout error: " + err.Error())
		os.Exit(1)
	}
	log.Printf("Api client initialized successfully with SID %v\n", apiClient.GetSessionID())
}
