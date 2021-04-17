package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type SSHAttempt struct {
	Host string `json:"host"`
	Log  string `json:"log"`
}

type SSHAttemptRecord struct {
	Id   string `json:"ID"`
	Host string `json:"host"`
	Log  string `json:"log"`
}

func main() {
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
	serverToken := fmt.Sprintf("Bearer %s", viper.GetString("AlphaServerToken"))

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello there stranger!")
	})

	app.Post("/ssh-attempts", func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		if bearer != serverToken {
			log.Warn(fmt.Sprintf("Invalid token from client: %s", bearer))
			return c.SendStatus(403)
		}

		log.Info(fmt.Sprintf("Received data from client: %s", c.Body()))
		event := SSHAttempt{}
		json.Unmarshal(c.Body(), &event)
		hash := sha256.New()
		hash.Write(c.Body())
		dbRecord := SSHAttemptRecord{
			Id:   hex.EncodeToString(hash.Sum(nil)),
			Host: event.Host,
			Log:  event.Log,
		}

		data, err := json.Marshal(dbRecord)
		if err != nil {
			panic(err)
		}
		log.Info(fmt.Sprintf("Store the record to db: %s", data))
		return c.Send(c.Body())
	})

	app.Listen(":3000")
}
