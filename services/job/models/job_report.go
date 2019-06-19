package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	// JobReportStatusPassed indicates that every no conflicts has been detected
	JobReportStatusPassed = "PASSED"
	// JobReportStatusWarning indicates that at least one conflict has been detected
	JobReportStatusWarning = "WARNING"
	// JobReportStatusError indicates that the job could not been executed due to an error (error details are store in the Message field)
	JobReportStatusError = "ERROR"
)

// JobReport is the data model used to represents a job report in the database
type JobReport struct {
	ID           uint `gorm:"primary_key"`
	CreationDate time.Time
	Message      string `gorm:"type:longtext"`
	Results      []*JobResult
	Status       string

	// Database references
	JobID uint `gorm:"index"`
}

// ReportStatistics is the datamodel used to represents report's statistics
type ReportStatistics struct {
	ResolvedCount uint
	NotFoundCount uint
	FoundCount    uint
}

// JobReportAccessor hosts every job report related methods
type JobReportAccessor struct {
	*gorm.DB
}

// NewJobReportAccessor helps to instanciate a new job report accessor with the given database connection
func NewJobReportAccessor(db *gorm.DB) *JobReportAccessor {
	return &JobReportAccessor{db}
}

// BeforeCreate helps to set the current date as creation date
func (report *JobReport) BeforeCreate() (err error) {
	report.CreationDate = time.Now()
	return nil
}

// Create add a new job report in the database
func (accessor *JobReportAccessor) Create(report *JobReport) error {
	return accessor.DB.Create(report).Error
}

// Save helps to persist modifications on a job report
func (accessor *JobReportAccessor) Save(report *JobReport) error {
	return accessor.DB.Save(report).Error
}

// GetByID returns a single job report from its ID
func (accessor *JobReportAccessor) GetByID(id uint) (*JobReport, error) {
	report := JobReport{}
	err := accessor.DB.First(&report, id).Error
	return &report, err
}

// GetByJobID returns every reports created by the given job
func (accessor *JobReportAccessor) GetByJobID(id uint) ([]*JobReport, error) {
	var reports []*JobReport
	err := accessor.DB.Where("job_id = ?", id).Find(&reports).Error
	return reports, err
}

// StatisticsByID returns statistics for the given report id
func (accessor *JobReportAccessor) StatisticsByID(id uint) (*ReportStatistics, error) {
	stats := &ReportStatistics{}
	var (
		resolved uint
		notFound uint
		found    uint
		err      error
	)

	err = accessor.DB.
		Table("job_results").
		Where("job_report_id = ? AND state LIKE 'RESOLVED'", id).
		Where("").
		Count(&resolved).
		Error
	if err != nil {
		return stats, err
	}
	stats.ResolvedCount = resolved

	err = accessor.DB.
		Table("job_results").
		Where("job_report_id = ? AND state LIKE 'FOUND'", id).
		Count(&found).
		Error
	if err != nil {
		return stats, err
	}
	stats.FoundCount = found

	err = accessor.DB.
		Table("job_results").
		Where("job_report_id = ? AND state LIKE 'NOT_FOUND'", id).
		Count(&notFound).
		Error
	if err != nil {
		return stats, err
	}
	stats.NotFoundCount = notFound

	return stats, err
}
