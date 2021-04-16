package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"net/http"

	"github.com/nxadm/tail"
	log "github.com/sirupsen/logrus"
)

var logFile = "/var/log/secure"
var alphaServerEndpoint = "https://localhost:8989/events"
var alphaServerToken = "blablabla"

func main() {
	t, err := tail.TailFile(logFile, tail.Config{Follow: true})
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
		}
	}
}

func isSSHAttempt(line string) bool {
	// Assume ssh attempt log contains the following line: [time] [pid]: ssh attempt from ip [ip address] using [user]
	var sshAttemptPattern = regexp.MustCompile(`.*: ssh attempt from ip .* using .*`)
	if sshAttemptPattern.MatchString(line) {
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

	req, err := http.NewRequest("POST", alphaServerEndpoint, strings.NewReader(line))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", alphaServerToken))
	_, err = client.Do(req)

	if err != nil {
		return err
	}
	return nil
}
