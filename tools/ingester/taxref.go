package ingester

import (
	"fmt"
	"strings"

	"git.ve.home/nicolasc/linotte/libs/csv"
	"git.ve.home/nicolasc/linotte/libs/helpers"
	"git.ve.home/nicolasc/linotte/models"
	"git.ve.home/nicolasc/linotte/tools/ingester/entries"
)

type TaxrefIngester struct {
	taxonAccessor          *models.TaxonAccessor
	verbAccessor           *models.VerbAccessor
	rankAccessor           *models.RankAccessor
	vernacularNameAccessor *models.VernacularNameAccessor
}

var taxrefIngester *TaxrefIngester = nil

func InitializeTaxrefIngester(
	taxonAccessor *models.TaxonAccessor,
	verbAccessor *models.VerbAccessor,
	rankAccessor *models.RankAccessor,
	vernacularNameAccessor *models.VernacularNameAccessor) *TaxrefIngester {

	if taxrefIngester == nil {
		taxrefIngester = &TaxrefIngester{
			taxonAccessor,
			verbAccessor,
			rankAccessor,
			vernacularNameAccessor,
		}
	}

	return taxrefIngester
}

func (ingester *TaxrefIngester) storeRanks(entries []entries.TaxrefEntry) map[string]uint {
	ranks := make(map[string]int)
	ids := make(map[string]uint)

	for _, element := range entries {
		ranks[element.RANG] = ranks[element.RANG] + 1
	}

	for id := range ranks {
		rank := models.Rank{Key: id}

		_, err := ingester.rankAccessor.Create(&rank)
		helpers.HandleError(err)

		ids[rank.Key] = rank.ID
	}

	return ids
}

func splitVernacularNames(raw string) []string {
	elements := strings.Split(raw, ",")
	var result []string

	for _, e := range elements {
		e = strings.Trim(e, " ")
		if e != "" && e != " " {
			result = append(result, e)
		}
	}

	return result
}

func (ingester *TaxrefIngester) storeVernacularNames(entries []entries.TaxrefEntry) map[string]models.VernacularName {
	names := make(map[string]int)
	results := make(map[string]models.VernacularName)

	for _, element := range entries {
		ns := splitVernacularNames(element.NOM_VERN)
		for _, n := range ns {
			names[n] = names[n] + 1
		}
	}

	for name := range names {
		vernacularName := models.VernacularName{Value: name}

		_, err := ingester.vernacularNameAccessor.Create(&vernacularName)
		helpers.HandleError(err)

		results[vernacularName.Value] = vernacularName
	}

	return results
}

func getVernacularNameEntities(raw string, list map[string]models.VernacularName) []models.VernacularName {
	var entities []models.VernacularName
	names := splitVernacularNames(raw)

	for _, name := range names {
		entities = append(entities, list[name])
	}

	return entities
}

func (ingester *TaxrefIngester) storeVerbsAndTaxons(entries []entries.TaxrefEntry, ranks map[string]uint, vernacularNamesList map[string]models.VernacularName) {
	for _, element := range entries {
		ingester.storeVerb(element, vernacularNamesList)

		_, err := ingester.taxonAccessor.Create(&models.Taxon{
			ID:               element.CD_NOM,
			ReferenceTaxonID: element.CD_REF,
			ParentID:         element.CD_SUP,
			RankID:           ranks[element.RANG],
		})
		helpers.HandleError(err)
	}
}

func (ingester *TaxrefIngester) storeVerb(entry entries.TaxrefEntry, vernacularNamesList map[string]models.VernacularName) {
	_, err := ingester.verbAccessor.Create(&models.Verb{
		VernacularNames: getVernacularNameEntities(entry.NOM_VERN, vernacularNamesList),
		Name:            entry.LB_NOM,
		Author:          entry.LB_AUTEUR,
		FirstLevelVernacularGroup:  entry.GROUP1_INPN,
		SecondLevelVernacularGroup: entry.GROUP2_INPN,
		ReferenceTaxonID:           entry.CD_REF,
		TaxonID:                    entry.CD_NOM,
	})
	helpers.HandleError(err)
}

func (ingester *TaxrefIngester) Ingest(file string) error {
	var (
		elements []entries.TaxrefEntry
		c        chan csv.Line = make(chan csv.Line)
	)

	fmt.Print("Parsing Taxref ... ")
	go csv.ReadFile(file, c, ';')
	for msg := range c {
		elements = append(elements, entries.ParseTaxrefEntry(msg.Elements))
	}
	fmt.Printf("File parsed. %d rows found.\n", len(elements))

	fmt.Print("Storing ranks ... ")
	ranks := ingester.storeRanks(elements)
	fmt.Printf("Ranks stored. %d rows inserted.\n", len(ranks))

	fmt.Print("Storing vernacular names ... ")
	names := ingester.storeVernacularNames(elements)
	fmt.Printf("Vernacular names stored. %d rows inserted.\n", len(names))

	fmt.Print("Storing verbs & taxons ... ")
	ingester.storeVerbsAndTaxons(elements, ranks, names)
	fmt.Println("Verbs inserted.")

	return nil
}
