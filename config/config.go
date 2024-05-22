package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func init() {
	var err error

	configFilePath := getConfigFilePath()

	config, err = loadConfig(configFilePath)
	if err != nil {
		log.Fatalln("Error loading config:", err)
	}
}

var config *Configuration

type Configuration struct {
	Telegram struct {
		PhoneNumber string `json:"phone_number"`
		APIID       int    `json:"api_id"`
		APIHash     string `json:"api_hash"`
	} `json:"telegram"`
	ClaimTo          string   `json:"claim_to"`
	CheckSleepTimeMS int      `json:"sleep_between_check"`
	Usernames        []string `json:"usernames"`
}

func GetConfig() *Configuration {
	return config
}

func loadConfig(path string) (*Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	config := &Configuration{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

func getConfigFilePath() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalln("Error getting wd:", err)
	}

	if filepath.Base(currentDir) == "app" {
		return filepath.Join(currentDir, "..", "..", "config.json")
	}

	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalln("Error getting executable path:", err)
	}

	executableDir := filepath.Dir(executablePath)
	return filepath.Join(executableDir, "config.json")
}
