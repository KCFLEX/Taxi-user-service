package config

import "github.com/spf13/viper"

type Config struct {
	Port      string `mapstructure:"HTTP_PORT"`
	Host      string `mapstructure:"HOST"`
	DbConn    string `mapstructure:"CONN_STRING"`
	SecretKey string `mapstructure:"JWT_SECRET_KEY"`
	RedisConn string `mapstructure:"REDIS_CONN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
