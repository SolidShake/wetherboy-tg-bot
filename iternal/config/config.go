package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	Bot struct {
		Token string `yaml:"token"`
	} `yaml:"bot"`

	OpenWether struct {
		Token string `yaml:"apptoken"`
	} `yaml:"openweather"`

	MongoDb struct {
		Host     string `yaml:"mongoHost"`
		Port     string `yaml:"mongoPort"`
		Database string `yaml:"mongoDatabase"`
	} `yaml:"MongoProperties"`
}

// GetConfig load configuration from yaml
func GetConfig() Config {
	var configFile *os.File
	if _, err := os.Stat("config.local.yaml"); os.IsNotExist(err) {
		configFile, _ = os.Open("config.yaml")
	} else {
		configFile, err = os.Open("config.local.yaml")
		if err != nil {
			panic(err)
		}
	}
	defer configFile.Close()

	var cnf Config
	decoder := yaml.NewDecoder(configFile)
	err := decoder.Decode(&cnf)
	if err != nil {
		panic(err)
	}

	log.Println("Config load successful")

	return cnf
}
