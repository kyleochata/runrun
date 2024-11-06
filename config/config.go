package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func InitConfig(fileName string) *viper.Viper {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	config := viper.New()
	config.SetConfigName(fileName)
	config.AddConfigPath(".")
	config.AddConfigPath("$HOME")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("Error while parsing configuration file", err)
	}
	replaceEnvVariables(config)
	return config
}

func replaceEnvVariables(config *viper.Viper) {
	for _, key := range config.AllKeys() {
		value := config.GetString(key)
		if strings.Contains(value, "${") {
			newValue := os.ExpandEnv(value)
			config.Set(key, newValue)
		}
	}
}
