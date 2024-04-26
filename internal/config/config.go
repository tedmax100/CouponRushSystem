package config

import (
	"context"
	"errors"

	"github.com/spf13/viper"
	"github.com/tedmax100/CouponRushSystem/internal/log"
	"go.uber.org/zap"
)

type Config struct {
	DB    string `mapstructure:"DATABASE_URL"`
	Redis string `mapstructure:"REDIS_URL"`
	Port  string `mapstructure:"PORT"`
	MQ    string `mapstructure:"MQ_URL"`
}

// GetConfig TODO May need to be private
func GetConfig() *Config {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath("./")
	err := v.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Error(context.Background(), err, zap.String("msg", "The error is of type ConfigFileNotFoundError"))
		}
	}

	v.AutomaticEnv()
	var envList = [][]string{
		{"REDIS_URL"},
		{"DATABASE_URL"},
		{"PORT"},
		{"MQ_URL"},
	}
	for _, envOptions := range envList {
		err := v.BindEnv(envOptions...)
		if err != nil {
			log.Fatal(context.Background(), err)
		}
	}

	var C Config
	err = v.Unmarshal(&C)
	if err != nil {
		panic(err)
	}
	return &C
}
