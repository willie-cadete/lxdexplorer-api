package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Interval  int
	LogLevel  string
	LXD       LXD
	HostNodes []string
	Server    API
	MongoDB   MongoDB
}

type API struct {
	Bind string
	Port int
}

type LXD struct {
	TLSCertificate    string
	TLSKey            string
	CertificateVerify bool
}

type MongoDB struct {
	URI string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/lxd-explorear-api/")
	viper.AddConfigPath("$HOME/.lxd-explorear-api")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found; using defaults")
		}

		if _, ok := err.(viper.ConfigParseError); ok {
			log.Fatalf("Unable to parse config file, %v", err)
		}
	}

	var config *Config
	log.Println("Using config file:", viper.ConfigFileUsed())
	log.Println("Using settings:", viper.AllSettings())
	err := viper.Unmarshal(&config)
	return config, err

}
