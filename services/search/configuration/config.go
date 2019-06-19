package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Configuration contains configurable properties
type Configuration struct {
	ServiceEndpoint string `split_words:"true"`
	ElasticHost     string `split_words:"true"`
	ElasticPort     int    `split_words:"true" default:"9200"`
	ElasticIndex    string `split_words:"true"`
}

// Get returns the configuration readed from environment variables
func Get() *Configuration {
	var config Configuration
	err := envconfig.Process("linotte_search", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &config
}
