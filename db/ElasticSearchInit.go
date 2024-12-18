package db

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
)

var ElasticClient *elasticsearch.Client

func InitElasticSearchClient() {
	var err error
	caCertPath := "./http_ca.crt"
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("Error reading CA cert: %s", err)
	}
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://192.168.31.75:9200",
		},
		Username: "elastic",              // 如果需要
		Password: "MH2c9MNGKi0YJnpysUX7", // 如果需要
		CACert:   caCert,
	}
	ElasticClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
}
