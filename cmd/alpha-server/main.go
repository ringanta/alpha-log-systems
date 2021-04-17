package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	viper.SetDefault("DBPort", "5432")
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

	db, err := initDB()
	if err != nil {
		panic(fmt.Sprintf("Can't connect to db: %s", err))
	}
	db.AutoMigrate(&SSHAttemptRecord{})
	log.Info("Initialization of connection to the database is success")

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		// SELECT COUNT(ID), Host FROM ssh_attempt_records GROUP BY Host;
		//rows, _ := db.Raw("SELECT host, COUNT(id) FROM ssh_attempt_records GROUP BY host").Rows()
		//defer rows.Close()

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

		db.Create(&dbRecord)
		log.Info(fmt.Sprintf("Store the record to db: %s", data))
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

func initDB() (*gorm.DB, error) {
	dbHost := viper.GetString("DBHost")
	dbPort := viper.GetString("DBPort")
	dbUser := viper.GetString("DBUser")
	dbPassword := viper.GetString("DBPassword")
	dbName := viper.GetString("DBName")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
