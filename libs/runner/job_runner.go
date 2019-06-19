package runner

import (
	"fmt"

	"git.ve.home/nicolasc/linotte/models"
	"github.com/fatih/color"
)

type JobConfiguration struct {
	BasePath      string
	FileExtension string
}

type JobRunner struct {
	configuration JobConfiguration

	jobAccessor       *models.JobAccessor
	taxonAccessor     *models.TaxonAccessor
	jobReportAccessor *models.JobReportAccessor
}

func NewJobRunner(configuration JobConfiguration) (*JobRunner, error) {
	var (
		err        error
		connection *models.Connection
	)

	models.InitDatabase()
	if connection, err = models.Connect(); err != nil {
		return nil, err
	}

	jobAccessor := models.NewJobAccessor(connection.DB)
	taxonAccessor := models.NewTaxonAccessor(connection.DB)
	jobReportAccessor := models.NewJobReportAccessor(connection.DB)

	return &JobRunner{
		configuration,
		jobAccessor,
		taxonAccessor,
		jobReportAccessor,
	}, nil
}

func (runner *JobRunner) Run() {
	jobs := runner.jobAccessor.GetByStatus(models.JobStatusNew)

	for _, job := range jobs {
		color.Magenta("Running %s ...", job.Name)

		job.Status = models.JobStatusPending
		runner.jobAccessor.Save(&job)

		switch job.Type {
		case models.JobTypeRedList:
			job.Reports = append(job.Reports, runner.RunRedList(job.File))
		}

		job.Status = models.JobStatusIdle
		runner.jobAccessor.Save(&job)
		color.Green("Job saved: %s", job.Name)
		fmt.Printf("\n")
	}
}
