package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"git.ve.home/nicolasc/linotte/libs/csv"
	"git.ve.home/nicolasc/linotte/libs/helpers"
	"git.ve.home/nicolasc/linotte/models"
)

type entry struct {
	CD_NOM   uint
	NAME_LAT string
	NAME_FR  string
	CAT_LRR  string
	NAME_POP string
	COM      string
}

var accessor *models.TaxonAccessor

func (runner JobRunner) RunRedList(file string) models.JobReport {
	var (
		c       chan csv.Line = make(chan csv.Line)
		entries []entry
		report  models.JobReport
	)

	accessor = runner.taxonAccessor
	go csv.ReadFile(runner.getFullPath(file), c, ';')
	for msg := range c {
		if msg.Error != nil {
			report.Status = models.JobReportStatusError
			report.Message = fmt.Sprintf("Read error (%v)", msg.Error)
			return report
		}

		elem, err := parse(msg.Elements)
		if err != nil {
			report.Status = models.JobReportStatusError
			report.Message = fmt.Sprintf("Parse error (%v)", err)
			return report
		}

		entries = append(entries, elem)
	}

	if len(entries) == 0 {
		report.Status = models.JobReportStatusError
		report.Message = fmt.Sprintf("No entries for list")
	} else {
		err := process(entries, &report)
		if err != nil {
			report.Status = models.JobReportStatusError
			report.Message = fmt.Sprintf("Processing error (%v)", err)
		}
	}

	return report
}

func process(entries []entry, report *models.JobReport) error {
	report.Status = models.JobReportStatusPassed
	for _, e := range entries {
		id := findTaxonID(e)

		serialized, err := json.Marshal(e)
		if err != nil {
			return err
		}

		result := models.JobResult{
			SearchData: string(serialized[:]),
			Value:      e.CAT_LRR,
		}

		if id > 0 {
			result.State = models.JobResultFound
			result.TaxonID = id
		} else {
			result.State = models.JobResultNotFound
			report.Status = models.JobReportStatusWarning
		}

		report.Results = append(report.Results, result)
	}
	return nil
}

func findTaxonID(e entry) uint {
	var t = &models.Taxon{}

	if e.CD_NOM > 0 {
		t = accessor.GetReferenceByID(e.CD_NOM)
	}
	if t.ID == 0 && len(e.NAME_LAT) > 0 {
		t = accessor.GetReferenceByName(e.NAME_LAT)
		if t.ID == 0 {
			t = accessor.GetReferenceByFullName(e.NAME_LAT, true)
		}
	}
	if t.ID == 0 && len(e.NAME_FR) > 0 {
		t = accessor.GetReferenceByVernacularName(e.NAME_FR)
	}
	if t.ID == 0 {
		return lastChance(e.NAME_LAT).ID
	}
	return t.ID
}

func lastChance(name string) *models.Taxon {
	fmt.Printf("[NC] Last chance for: %s\n", name)
	var result []string

	name = strings.Replace(name, "L.", " ", -1)
	name = strings.Replace(name, " ET ", " ", -1)
	name = strings.Replace(name, " et ", " ", -1)
	name = strings.Replace(name, " & ", " ", -1)
	name = strings.Replace(name, " ex ", " ", -1)
	name = regexp.MustCompile(`\([^)]*\)`).ReplaceAllString(name, " ")
	name = regexp.MustCompile(`[0-9]+`).ReplaceAllString(name, " ")

	parts := strings.Split(name, " ")
	result = append(result, parts[0])
	for i, w := range parts {
		startsWithUpper, _ := regexp.MatchString(`^[A-Z].*`, w)
		if i > 0 && w != "" && w != " " && !startsWithUpper {
			result = append(result, w)
		}
	}

	term := strings.Join(result, " ")
	fmt.Printf("[NC] After all substitutions, we'll search on: %s\n", term)
	taxon := accessor.GetReferenceByName(term)
	if taxon.ID == 0 {
		taxon = accessor.GetReferenceByFullName(term, true)
	}

	if taxon.ID != 0 {
		fmt.Printf("[NC] Found !\n")
	} else {
		fmt.Printf("[NC] Not found ...\n")
	}

	return taxon
}

func parse(elements []string) (entry, error) {
	var (
		err error
		id  uint64
	)

	if len(elements) < 6 {
		return entry{}, errors.New("wrong columns count")
	}

	if elements[0] != "" {
		id, err = strconv.ParseUint(elements[0], 10, 64)
		helpers.HandleError(err)
	}

	return entry{
		CD_NOM:   uint(id),
		NAME_LAT: elements[1],
		NAME_FR:  elements[2],
		CAT_LRR:  elements[3],
		NAME_POP: elements[4],
		COM:      elements[5],
	}, nil
}

func (runner *JobRunner) getFullPath(file string) string {
	return fmt.Sprintf("%s%s.%s", runner.configuration.BasePath, file, runner.configuration.FileExtension)
}
