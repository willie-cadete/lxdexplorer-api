package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Interval  int
	Retention int32
	LogLevel  string
	LXD       LXD
	HostNodes []string
	Server    API
	MongoDB   MongoDB
}

type API struct {
	Bind string
	Port string
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

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found; using defaults")
			return nil, err
		}

		if _, ok := err.(viper.ConfigParseError); ok {
			log.Fatalf("Unable to parse config file, %v", err)
			return nil, err
		}
	}

	var config *Config
	log.Println("Using config file:", viper.ConfigFileUsed())

	for _, key := range viper.AllKeys() {
		envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		err := viper.BindEnv(key, envKey)
		if err != nil {
			log.Println("config: unable to bind env: " + err.Error())
		}
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return config, err

}
