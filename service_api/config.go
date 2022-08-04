package main

import (
	"fmt"

	"github.com/spf13/viper"
)

const config_file_name = "config"

func init_config() {
	viper.SetConfigName(config_file_name)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
