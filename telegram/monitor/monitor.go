package monitor

import (
	"app/telegram"
	"sync"
	"time"
)

func StartMonitor(usernames []string, sleepTimeMs int, availableUsernamesChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		usernames = checkUsernames(usernames, sleepTimeMs, availableUsernamesChan)
	}
}

func checkUsernames(usernames []string, sleepTimeMs int, availableUsernamesChan chan<- string) []string {
	for i := 0; i < len(usernames); i++ {
		username := usernames[i]
		if telegram.IsUsernameAvailable(username) {
			availableUsernamesChan <- username

			// Remove usernames that have already been passed to availableUsernamesChan
			usernames = append(usernames[:i], usernames[i+1:]...)
			i--
		}

		time.Sleep(time.Duration(sleepTimeMs) * time.Millisecond)
	}

	return usernames
}
