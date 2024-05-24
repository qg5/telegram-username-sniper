package monitor

import (
	"app/telegram"
	"sync"
	"time"
)

var mu sync.Mutex

func StartMonitor(usernames []string, sleepTimeMs int, availableUsernamesChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		checkUsernames(&usernames, sleepTimeMs, availableUsernamesChan)
	}
}

func checkUsernames(usernames *[]string, sleepTimeMs int, availableUsernamesChan chan<- string) {
	for i := 0; i < len(*usernames); i++ {
		username := (*usernames)[i]
		if telegram.IsUsernameAvailable(username) {
			availableUsernamesChan <- username

			mu.Lock()
			defer mu.Unlock()
			*usernames = append((*usernames)[:i], (*usernames)[i+1:]...)
		}

		time.Sleep(time.Duration(sleepTimeMs) * time.Millisecond)
	}
}
