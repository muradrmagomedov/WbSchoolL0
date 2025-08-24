package main

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("не получилось прочитать файл конфига")
	}
}
