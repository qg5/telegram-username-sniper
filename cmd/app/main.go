import (
	"app/config"
	"app/telegram"
	"app/telegram/monitor"
	"app/telegram/sniper"
	"fmt"
	"log"
	"strings"
	"sync"
)

var (
	cfg = config.GetConfig()
	wg  sync.WaitGroup
)

func main() {
	phoneNumber := cfg.Telegram.PhoneNumber
	if !strings.HasPrefix(phoneNumber, "+") {
		log.Fatal("Phone number must be in international format (e.g, +12848428429)")
	}

	claimMethod := cfg.ClaimTo
	if claimMethod != "user" && claimMethod != "channel" {
		log.Fatal("The value of 'claim_to' must be either 'user' or 'channel'")
	}

	usernames := cfg.Usernames
	if len(usernames) == 0 {
		log.Fatal("Please provide at least 1 username")
	}

	client := telegram.New(cfg.Telegram.APIID, cfg.Telegram.APIHash, phoneNumber)

	availableUsernamesChan := make(chan string)
	defer close(availableUsernamesChan)

	fmt.Printf("Monitoring %v username(s)\n", len(usernames))
	go monitor.StartMonitor(usernames, cfg.CheckSleepTimeMS, availableUsernamesChan, &wg)
	go sniper.ProcessAvailableUsernames(client, claimMethod, availableUsernamesChan)

	select {}
}
