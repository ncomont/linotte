package services

import (
	"context"
	"strconv"

	proto_job "git.ve.home/nicolasc/linotte/services/search/proto"
	elastic "gopkg.in/olivere/elastic.v5"
)

// SearchService give an access to every possible search operations
type SearchService struct {
	client *elastic.Client
	index  string
}

// NewSearchService creates a new search service from the given client
func NewSearchService(client *elastic.Client, index string) *SearchService {
	return &SearchService{
		client: client,
		index:  index,
	}
}

// Search returns every taxon ids corresponding to the given search parameters
func (service *SearchService) Search(req *proto_job.SearchRequest) (*proto_job.SearchReply, error) {
	var (
		ids         []uint32
		hasCriteria bool
		sortOn      = "id"
	)

	query := elastic.NewBoolQuery()

	if len(req.Term) > 0 {
		hasCriteria = true
		query = query.Must(elastic.NewMultiMatchQuery(req.Term, "name", "vernacularNames", "author").Fuzziness("AUTO"))
	}
	if len(req.Ranks) > 0 {
		hasCriteria = true
		ranksQuery := elastic.NewBoolQuery()
		for _, r := range req.Ranks {
			ranksQuery = ranksQuery.Should(elastic.NewMatchQuery("rank", r))
		}
		query = query.Must(ranksQuery)
	}
	if req.Reference {
		hasCriteria = true
		query = query.Must(elastic.NewMatchQuery("isReference", "true"))
	}

	if hasCriteria {
		sortOn = "_score"
	}

	builder := service.client.
		Search().
		Index(service.index).
		StoredFields("_id").
		Sort(sortOn, false).
		Query(query).
		From(int(req.Offset)).
		Size(int(req.Limit))

	searchResult, err := builder.Do(context.Background())
	if err != nil {
		return &proto_job.SearchReply{Total: 0}, err
	}

	for _, item := range searchResult.Hits.Hits {
		v, err := strconv.ParseUint(item.Id, 10, 64)
		if err != nil {
			return &proto_job.SearchReply{Total: 0}, err
		}
		ids = append(ids, uint32(v))
	}

	return &proto_job.SearchReply{
		Total:   searchResult.TotalHits(),
		Results: ids,
	}, nil
}
