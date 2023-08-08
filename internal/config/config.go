package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GetConfig()
}

var (
	c *Config
)

type Config struct {
	Mongo Mongo
}

type Mongo struct {
	Address string
	DBName  string
}

func GetConfig() *Config {
	if c == nil {
		//MONGO
		address := os.Getenv("MONGO_ADDRESS")
		if address == "" {
			panic("MONGO_ADDRESS is not set")
		}

		name := os.Getenv("MONGO_DB_NAME")
		if name == "" {
			panic("MONGO_DB_NAME is not set")
		}

		c = &Config{
			Mongo: Mongo{
				Address: address,
				DBName:  name,
			},
		}

		return c
	}

	return c
}
