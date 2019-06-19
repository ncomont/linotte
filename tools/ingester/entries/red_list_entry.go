package entries

import (
	"strconv"

	"git.ve.home/nicolasc/linotte/libs/helpers"
	"git.ve.home/nicolasc/linotte/models"
)

type RedListEntry struct {
	NAME_LR     string
	NAME_LOC    string
	DATE        int
	AUTHORS     string
	REF_LRR     string
	NAME_FOLDER string
	Entity      models.RedList
}

func ParseRedListEntry(record []string) RedListEntry {
	var (
		date int
		err  error
	)

	if len(record[4]) > 0 {
		date, err = strconv.Atoi(record[4])
		helpers.HandleError(err)
	}

	return RedListEntry{
		NAME_LR:     record[1],
		NAME_LOC:    record[3],
		DATE:        date,
		AUTHORS:     record[5],
		REF_LRR:     record[6],
		NAME_FOLDER: record[7],
	}
}
