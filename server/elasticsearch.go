package main

import (
	"context"
	"time"

	"github.com/knarfeh/searchtelegram/server/diagnose"
	elastic "gopkg.in/olivere/elastic.v5"
)

// ElasticConfig is configuration for Elasticsearch instance
type ElasticConfig struct {
	Endpoint           string
	Username           string
	Password           string
	Retries            int
	HealthCheckTimeout time.Duration
}

// GetClientOption returns a elasticSearch client option
func (c ElasticConfig) GetClientOption() []elastic.ClientOptionFunc {
	options := make([]elastic.ClientOptionFunc, 4, 5)
	options[0] = elastic.SetURL(c.Endpoint)
	options[1] = elastic.SetMaxRetries(c.Retries)
	options[2] = elastic.SetSniff(false)
	options[3] = elastic.SetHealthcheckTimeoutStartup(c.HealthCheckTimeout)
	if c.Username != "" {
		options = append(options, elastic.SetBasicAuth(c.Username, c.Password))
	}
	return options
}

// ESClient ...
type ESClient struct {
	Client *elastic.Client
	config ElasticConfig
}

// NewESClient ...
func NewESClient(config ElasticConfig) (*ESClient, error) {

	client, err := elastic.NewClient(
		config.GetClientOption()...,
	)

	return &ESClient{
		Client: client,
		config: config,
	}, err
}

// Diagnose runs a diagnose over ES connection
func (es *ESClient) Diagnose() diagnose.ComponentReport {
	return diagnose.SimpleDiagnose("elastic_search", func() error {
		_, _, err := es.Client.Ping(es.config.Endpoint).Do(context.TODO())
		return err
	})
}
