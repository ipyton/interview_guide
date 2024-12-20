package db

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

var ElasticClient *elasticsearch.Client

func InitElasticSearchClient() {
	var err error
	//caCertPath := "./http_ca.crt"
	//caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("Error reading CA cert: %s", err)
	}
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://113.44.65.250:9200",
		},
		Username: "elastic",              // 如果需要
		Password: "Om5-S+1vLIXyso5Zx3*2", // 如果需要

	}
	ElasticClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
}
