package logic

import (
	"context"
	"ngoc/domain"

	"github.com/elastic/go-elasticsearch/v8"
)

func ConnectWithElasticsearch(ctx context.Context) context.Context {

	newClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	})
	if err != nil {
		panic(err)
	}

	return context.WithValue(ctx, domain.ClientKey, newClient)

}
