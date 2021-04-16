package main

import (
	"fmt"
	"regexp"
	"strings"
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

func main() {
	viper.SetDefault("SSHLogFile", "/var/log/auth.log")
	viper.SetDefault("AlphaServerEndpoint", "http://localhost:3000/")
	viper.SetDefault("AlphaServerToken", "")
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
			log.Info(fmt.Sprintf("Sending ssh attempt: %s", line.Text))
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

	req, err := http.NewRequest("POST", viper.GetString("AlphaServerEndpoint"), strings.NewReader(line))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", viper.GetString("AlphaServerToken")))
	_, err = client.Do(req)

	if err != nil {
		return err
	}
	return nil
}
