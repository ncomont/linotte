package ingester

import (
	"errors"

	"git.ve.home/nicolasc/linotte/models"
	"git.ve.home/nicolasc/linotte/settings"
)

var (
	instance *Ingester = nil
	path     string
)

type Ingester struct {
	Taxref       *TaxrefIngester
	RedList      *RedListIngester
	User         *UserIngester
	SearchEngine *SearchEngineIngester
}

func Initialize() (*Ingester, error) {
	var (
		err        error
		connection *models.Connection
	)

	if instance == nil {
		models.InitDatabase()
		if connection, err = models.Connect(); err == nil {
			if err = models.MigrateAll(); err == nil {
				instance = &Ingester{
					Taxref: InitializeTaxrefIngester(
						models.NewTaxonAccessor(connection.DB),
						models.NewVerbAccessor(connection.DB),
						models.NewRankAccessor(connection.DB),
						models.NewVernacularNameAccessor(connection.DB),
					),
					RedList: InitializeRedListIngester(
						models.NewRedListAccessor(connection.DB),
						models.NewJobAccessor(connection.DB),
					),
					User: InitializeUserIngester(
						models.NewUserAccessor(connection.DB),
					),
					SearchEngine: InitializeSearchEngineIngester(
						models.NewTaxonAccessor(connection.DB),
					),
				}

				path = settings.Get().VolumesPath
			}
		}
	}

	return instance, err
}

func (ingester *Ingester) StartIngest(dataset string) error {
	switch dataset {
	case "taxref":
		ingester.ingestTaxref()
		break
	case "redlist":
		ingester.ingestRedList()
		break
	case "user":
		ingester.ingestUser()
		break
	case "searchengine":
		ingester.ingestSearchEngine()
		break
	default:
		return errors.New("Unknown runner")
	}

	return nil
}

func (ingester *Ingester) ingestTaxref() error {
	return ingester.Taxref.Ingest(path + "TAXREFv10.0.csv")
}

func (ingester *Ingester) ingestRedList() error {
	return ingester.RedList.Ingest(path + "redlists/")
}

func (ingester *Ingester) ingestUser() error {
	return ingester.User.Ingest()
}

func (ingester *Ingester) ingestSearchEngine() error {
	return ingester.SearchEngine.Ingest()
}

func ingestZnieff() {}

func ingestStatus() {}
