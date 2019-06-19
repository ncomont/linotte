package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Configuration contains configurable properties
type Configuration struct {
	ServiceEndpoint     string `split_words:"true"`
	DatabaseHost        string `split_words:"true"`
	DatabaseUser        string `split_words:"true"`
	DatabasePassword    string `split_words:"true"`
	DatabaseName        string `split_words:"true"`
	DatabasePort        int    `split_words:"true" default:"3306"`
	DatabaseVerboseMode bool   `split_words:"true" default:"false"`
	PublicKeyPath       string `split_words:"true"`
	PrivateKeyPath      string `split_words:"true"`
}

// Get returns the configuration readed from environment variables
func Get() *Configuration {
	var config Configuration
	err := envconfig.Process("linotte_user", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &config
}
