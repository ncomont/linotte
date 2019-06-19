package search

import (
	"context"
	"fmt"
	"strconv"

	"git.ve.home/nicolasc/linotte/models"
	"git.ve.home/nicolasc/linotte/settings"
	"gopkg.in/olivere/elastic.v5"
)

var (
	instance *Engine
)

type IndexedTaxon struct {
	ID              uint     `json:"id"`
	Name            string   `json:"name"`
	VernacularNames []string `json:"vernacularNames"`
	Author          string   `json:"author"`
	Rank            string   `json:"rank"`
	IsReference     bool     `json:"isReference"`
}

type Engine struct {
	client *elastic.Client
}

func InitializeEngine() (*Engine, error) {
	var (
		err    error
		client *elastic.Client
	)

	if instance == nil {
		url := fmt.Sprintf("http://%s:%s", settings.Get().ElasticSearchHost, settings.Get().ElasticSearchPort)
		if client, err = elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(true)); err == nil {
			instance = &Engine{client}
		}
	}

	return instance, err
}

func (engine *Engine) StartIndexer(accessor *models.TaxonAccessor) error {
	var (
		err     error
		indexer *Indexer
	)

	if indexer, err = InitializeIndexer(accessor, engine.client); err == nil {
		indexer.Start(1000)
	}
	return err
}

type SP struct {
	Term          string
	ReferenceOnly bool
	Ranks         []string
}

type SearchResult struct {
	Total  int
	Values []uint
}

func (engine *Engine) Search(limit int, offset int, sp SP) (SearchResult, error) {
	var (
		ids         []uint
		result      SearchResult
		hasCriteria bool
		sortOn      = "id"
	)

	query := elastic.NewBoolQuery()

	if len(sp.Term) > 0 {
		hasCriteria = true
		query = query.Must(elastic.NewMultiMatchQuery(sp.Term, "name", "vernacularNames", "author").Fuzziness("AUTO"))
	}
	if len(sp.Ranks) > 0 {
		hasCriteria = true
		ranksQuery := elastic.NewBoolQuery()
		for _, r := range sp.Ranks {
			ranksQuery = ranksQuery.Should(elastic.NewMatchQuery("rank", r))
		}
		query = query.Must(ranksQuery)
	}
	if sp.ReferenceOnly {
		hasCriteria = true
		query = query.Must(elastic.NewMatchQuery("isReference", "true"))
	}

	if hasCriteria {
		sortOn = "_score"
	}

	builder := engine.client.
		Search().
		Index(settings.Get().ElasticSearchIndex).
		StoredFields("_id").
		Sort(sortOn, false).
		Query(query).
		From(offset).
		Size(limit)

	searchResult, err := builder.Do(context.Background())
	if err != nil {
		return result, err
	}

	for _, item := range searchResult.Hits.Hits {
		v, err := strconv.ParseUint(item.Id, 10, 64)
		if err != nil {
			return result, err
		}
		ids = append(ids, uint(v))
	}

	return SearchResult{
		Total:  int(searchResult.TotalHits()),
		Values: ids,
	}, nil
}
