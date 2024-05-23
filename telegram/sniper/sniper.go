package sniper

import (
	"app/telegram"
	"fmt"
	"log"
)

func ProcessAvailableUsernames(client *telegram.Client, claimMethod string, availableUsernamesChan <-chan string) {
	for {
		username := <-availableUsernamesChan
		fmt.Println("Found available username:", username)

		err := claimUsername(client, claimMethod, username)
		if err != nil {
			log.Println("Failed to claim:", username)
		}
	}
}

func claimUsername(client *telegram.Client, claimMethod, username string) error {
	var err error
	switch claimMethod {
	case "channel":
		err = client.CreateChannel(username)
	case "user":
		err = client.UpdateUsername(username)
	}

	return err
}
