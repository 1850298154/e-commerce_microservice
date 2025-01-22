package meili

import (
	"2501YTC/app/product/conf"

	"github.com/meilisearch/meilisearch-go"
)

var Client meilisearch.ServiceManager

func Init() {
	Client = meilisearch.New(conf.GetConf().Meili.Address, meilisearch.WithAPIKey(conf.GetConf().Meili.APIKey))
}
