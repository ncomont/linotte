package entries

import (
	"strconv"

	"git.ve.home/nicolasc/linotte/libs/helpers"
)

type TaxrefEntry struct {
	CD_NOM      uint
	CD_REF      uint
	CD_SUP      uint
	RANG        string
	NOM_COMPLET string
	GROUP1_INPN string
	GROUP2_INPN string
	NOM_VERN    string
	LB_NOM      string
	LB_AUTEUR   string
}

func ParseTaxrefEntry(record []string) TaxrefEntry {
	var nom, sup, ref uint64
	var err error

	if len(record[7]) > 0 {
		nom, err = strconv.ParseUint(record[7], 10, 64)
		helpers.HandleError(err)
	}
	if len(record[9]) > 0 {
		sup, err = strconv.ParseUint(record[9], 10, 64)
		helpers.HandleError(err)
	}
	if len(record[10]) > 0 {
		ref, err = strconv.ParseUint(record[10], 10, 64)
		helpers.HandleError(err)
	}

	return TaxrefEntry{
		CD_NOM:      uint(nom),
		CD_SUP:      uint(sup),
		CD_REF:      uint(ref),
		NOM_COMPLET: record[14],
		RANG:        record[11],
		GROUP1_INPN: record[5],
		GROUP2_INPN: record[6],
		NOM_VERN:    record[17],
		LB_NOM:      record[12],
		LB_AUTEUR:   record[13],
	}
}

func (e TaxrefEntry) IsReference() bool {
	return e.CD_NOM == e.CD_REF
}
