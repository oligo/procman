package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func LoadConfig() {

	viper.SetConfigName("apps")
	viper.AddConfigPath(".")
	//viper.AddConfigPath("")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error %s\n", err))
	}
}
