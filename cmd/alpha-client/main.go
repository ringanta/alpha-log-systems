package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"net/http"

	"github.com/nxadm/tail"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	sshAttemptSuccess     = regexp.MustCompile(`sshd\[.*\]: Accepted`)
	sshAttemptInvalidUser = regexp.MustCompile(`sshd\[.*\]: Invalid user`)
	sshAttemptInvalidCred = regexp.MustCompile(`sshd\[.*\]: Connection closed`)
)

type sshAttempt struct {
	Host string `json:"host"`
	Log  string `json:"log"`
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	viper.SetDefault("SSHLogFile", "/var/log/auth.log")
	viper.SetDefault("AlphaServerEndpoint", "http://localhost:3000/")
	viper.SetDefault("AlphaServerToken", "")
	viper.SetDefault("Hostname", hostname)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("Config file not found will proceed with default")
		} else {
			panic(fmt.Errorf("Fatal error when parsing config file: %s \n", err))
		}
	}

	t, err := tail.TailFile(viper.GetString("SSHLogFile"), tail.Config{Follow: true})
	if err != nil {
		panic(err)
	}

	for line := range t.Lines {
		if isSSHAttempt(line.Text) {
			err := sendEventToAlphaServer(line.Text)
			if err != nil {
				// Handle failure sending to server
				log.Error(fmt.Sprintf("Failed to send event to server: %s", err))
			}
		} else {
			log.Info(fmt.Sprintf("Not connection attempt, skipping: %s", line.Text))
		}
	}
}

func isSSHAttempt(line string) bool {
	if sshAttemptSuccess.MatchString(line) ||
		sshAttemptInvalidCred.MatchString(line) ||
		sshAttemptInvalidUser.MatchString(line) {
		return true
	}
	return false
}

func sendEventToAlphaServer(line string) error {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	payload, err := json.Marshal(&sshAttempt{
		Host: viper.GetString("Hostname"),
		Log:  line,
	})

	if err != nil {
		return fmt.Errorf("Failed to construct sshAttempt struct: %s", err)
	}

	req, err := http.NewRequest("POST", viper.GetString("AlphaServerEndpoint"), bytes.NewBuffer(payload))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", viper.GetString("AlphaServerToken")))
	log.Info(fmt.Sprintf("Sending ssh attempt: %s", payload))
	_, err = client.Do(req)

	if err != nil {
		return err
	}
	return nil
}
