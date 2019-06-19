package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Configuration contains configurable properties
type Configuration struct {
	RabbitEndpoint        string `split_words:"true"`
	RabbitTaskQueueID     string `split_words:"true"`
	RabbitResultQueueID   string `split_words:"true"`
	RabbitStatusQueueID   string `split_words:"true"`
	IndexerBatchSize      string `split_words:"true"`
	TaxrefServiceEndpoint string `split_words:"true"`
	StoragePath           string `split_words:"true"`
}

// Get returns the configuration readed from environment variables
func Get() *Configuration {
	var config Configuration
	err := envconfig.Process("linotte_job_runner", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &config
}
