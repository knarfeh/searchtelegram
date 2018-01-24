package main

import (
	"fmt"
	"strings"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

// ESClient ...
type ESClient struct {
	Client *elastic.Client
}

func (es *ESClient) Test() bool {
	return es.Client.IsRunning()
}

// NewESClient ...
func NewESClient(endpoint string, username string, password string, retries int) (*ESClient, error) {
	if !strings.Contains(endpoint, "http") {
		fmt.Println("Adding http to endpoint:", endpoint)
		endpoint = "http://" + endpoint
	}
	options := make([]elastic.ClientOptionFunc, 4, 5)
	options[0] = elastic.SetURL(endpoint)
	options[1] = elastic.SetMaxRetries(retries)
	options[2] = elastic.SetSniff(false)
	options[3] = elastic.SetHealthcheckTimeoutStartup(3 * time.Second)
	if username != "" {
		options = append(options, elastic.SetBasicAuth(username, password))
	}

	client, err := elastic.NewClient(
		options...,
	)

	return &ESClient{
		Client: client,
	}, err
}
