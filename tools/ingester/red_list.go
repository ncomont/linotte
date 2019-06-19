package ingester

import (
	"encoding/json"
	"fmt"
	"time"

	"git.ve.home/nicolasc/linotte/libs/csv"
	"git.ve.home/nicolasc/linotte/libs/helpers"
	"git.ve.home/nicolasc/linotte/models"
	"git.ve.home/nicolasc/linotte/tools/ingester/entries"
)

type RedListIngester struct {
	redListAccessor *models.RedListAccessor
	jobAccessor     *models.JobAccessor
}

var redListIngester *RedListIngester = nil

func InitializeRedListIngester(
	redListAccessor *models.RedListAccessor,
	jobAccessor *models.JobAccessor) *RedListIngester {

	if redListIngester == nil {
		redListIngester = &RedListIngester{
			redListAccessor,
			jobAccessor,
		}
	}

	return redListIngester
}

func getIndexPath(path string) string {
	return path + "index.csv"
}

func (ingester *RedListIngester) Ingest(path string) error {
	var (
		c     chan csv.Line = make(chan csv.Line)
		count int
	)

	fmt.Println("Parsing red list index ... ")
	go csv.ReadFile(getIndexPath(path), c, ';')
	for msg := range c {
		index := entries.ParseRedListEntry(msg.Elements)
		blob, err := json.Marshal(models.RedList{
			Name:         index.NAME_LR,
			OriginalName: index.REF_LRR,
			Author:       index.AUTHORS,
			Date:         index.DATE,
			// Area MUST be set here
		})
		helpers.HandleError(err)

		_, err = ingester.jobAccessor.Create(&models.Job{
			Name:       index.NAME_LR,
			File:       index.NAME_FOLDER,
			Status:     models.JobStatusIdle,
			Type:       models.JobTypeRedList,
			LastUpdate: time.Now(),
			Data:       string(blob[:]),
		})
		helpers.HandleError(err)

		// JobRunner.Run(job.ID) ?

		count++
	}
	fmt.Printf("File parsed. %d jobs added.\n", count)

	return nil
}
