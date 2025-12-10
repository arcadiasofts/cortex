package config

import (
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	App struct {
		Name    string `mapstructure:"name"`
		Port    string `mapstructure:"port"`
		Version string `mapstructure:"version"`
	} `mapstructure:"app"`
	DB struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Pass     string `mapstructure:"pass"`
		Name     string `mapstructure:"name"`
		TimeZone string `mapstructure:"timezone"`
	} `mapstructure:"db"`
	Redis struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}
	Jwt struct {
		Secret string `mapstructure:"secret"`
	}
}

func ProvideConfig() (*AppConfig, error) {
	v := viper.New()

	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Config environment variables override settings (For Docker/K8s)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config files
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var c AppConfig
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
