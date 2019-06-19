package ingester

import (
	"fmt"

	"git.ve.home/nicolasc/linotte/libs/search"
	"git.ve.home/nicolasc/linotte/models"
)

type SearchEngineIngester struct {
	taxonAccessor *models.TaxonAccessor
}

var searchEngineIngester *SearchEngineIngester = nil

func InitializeSearchEngineIngester(
	taxonAccessor *models.TaxonAccessor) *SearchEngineIngester {

	if searchEngineIngester == nil {
		searchEngineIngester = &SearchEngineIngester{taxonAccessor}
	}

	return searchEngineIngester
}

func (ingester *SearchEngineIngester) Ingest() error {
	fmt.Println("Indexing taxons ... ")

	engine, err := search.InitializeEngine()
	if err != nil {
		return err
	}

	err = engine.StartIndexer(ingester.taxonAccessor)
	if err != nil {
		return err
	}

	fmt.Println("Taxons indexed.")

	return nil
}
