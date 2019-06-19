package search

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"sync/atomic"
	"time"

	"git.ve.home/nicolasc/linotte/models"
	"git.ve.home/nicolasc/linotte/settings"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"gopkg.in/olivere/elastic.v5"
)

var (
	indexerInstance *Indexer
	ctx             = context.Background()
)

type DataPayload struct {
	ID    string
	Taxon *IndexedTaxon
}

type Indexer struct {
	accessor *models.TaxonAccessor
	client   *elastic.Client
}

func InitializeIndexer(accessor *models.TaxonAccessor, client *elastic.Client) (*Indexer, error) {
	if indexerInstance == nil {
		indexerInstance = &Indexer{
			accessor: accessor,
			client:   client,
		}
	}

	return indexerInstance, nil
}

func (indexer *Indexer) Start(bulkSize int) error {
	docsc := make(chan *DataPayload)
	g, ctx := errgroup.WithContext(context.TODO())
	begin := time.Now()

	g.Go(func() error {
		defer close(docsc)
		size := settings.Get().IndexerBatchSize
		count := float64(indexer.accessor.GetTotalCount())
		pages := int(math.Ceil(float64(count) / float64(size)))
		for i := 0; i < pages; i++ {
			taxons := indexer.accessor.GetAll(size, i*size)

			for _, t := range taxons {
				var vn []string
				for _, f := range t.Verb.VernacularNames {
					vn = append(vn, f.Value)
				}

				pl := &DataPayload{
					ID: strconv.Itoa(int(t.ID)),
					Taxon: &IndexedTaxon{
						ID:              t.ID,
						Name:            t.Verb.Name,
						VernacularNames: vn,
						Author:          t.Verb.Author,
						Rank:            t.Rank.Key,
						IsReference:     t.ID == t.ReferenceTaxonID,
					},
				}

				select {
				case docsc <- pl:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}

		return nil
	})

	var total uint64
	g.Go(func() error {
		index := settings.Get().ElasticSearchIndex

		if err := clearIndex(ctx, indexer.client, index); err != nil {
			return err
		}

		bulk := indexer.client.Bulk().Index(index).Type("verb")
		for d := range docsc {
			// Simple progress
			current := atomic.AddUint64(&total, 1)
			dur := time.Since(begin).Seconds()
			sec := int(dur)
			pps := int64(float64(current) / dur)
			fmt.Printf("%10d | %6d req/s | %02d:%02d\r", current, pps, sec/60, sec%60)

			bulk.Add(elastic.NewBulkIndexRequest().Id(d.ID).Doc(d.Taxon))
			if bulk.NumberOfActions() >= bulkSize {
				res, err := bulk.Do(ctx)
				if err != nil {
					return err
				}
				if res.Errors {
					return errors.New("bulk commit failed")
				}
			}

			select {
			default:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if bulk.NumberOfActions() > 0 {
			if _, err := bulk.Do(ctx); err != nil {
				return err
			}
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	dur := time.Since(begin).Seconds()
	sec := int(dur)
	pps := int64(float64(total) / dur)
	fmt.Printf("%10d | %6d req/s | %02d:%02d\n", total, pps, sec/60, sec%60)

	return nil
}

func clearIndex(ctx context.Context, client *elastic.Client, index string) error {
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return err
	}
	if exists {
		_, err = client.DeleteIndex(index).Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
