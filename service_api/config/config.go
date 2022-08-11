package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const config_file_name = "config"

func InitConfig() {
	viper.SetConfigName(config_file_name)
	viper.AddConfigPath(".")
	viper.BindEnv("POSTGRES_USER")
	viper.BindEnv("POSTGRES_PASSWORD")
	viper.BindEnv("POSTGRES_HOST")
	viper.BindEnv("POSTGRES_PORT")
	viper.BindEnv("POSTGRES_DB")
	viper.BindEnv("KAFKA_URL")
	viper.BindEnv("OBJECT_CREATION_TOPIC_NAME")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func GetDbConnectionString() string {
	user := viper.GetString("POSTGRES_USER")
	pwd := viper.GetString("POSTGRES_PASSWORD")
	host := viper.GetString("POSTGRES_HOST")
	port := viper.GetInt("POSTGRES_PORT")
	db := viper.GetString("POSTGRES_DB")
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", user, pwd, host, port, db)
}
