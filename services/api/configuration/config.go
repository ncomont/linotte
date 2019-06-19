package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Configuration contains configurable properties
type Configuration struct {
	ApiServiceEndpoint    string `split_words:"true"`
	TaxrefServiceEndpoint string `split_words:"true"`
	JobServiceEndpoint    string `split_words:"true"`
	UserServiceEndpoint   string `split_words:"true"`
	SearchServiceEndpoint string `split_words:"true"`
}

// Get returns the configuration readed from environment variables
func Get() *Configuration {
	var config Configuration
	err := envconfig.Process("linotte", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &config
}
