package models

import "github.com/jinzhu/gorm"

const (
	// JobResultFound inidicates that a correspondance has been found
	JobResultFound = 0
	// JobResultNotFound indicates that no correspondance has been found
	JobResultNotFound = 1
	// JobResultResolved indicateds that the conflict has been manually resolved
	JobResultResolved = 2
)

// JobResult represent the data model for a job result
type JobResult struct {
	ID         uint `gorm:"primary_key"`
	State      string
	SearchData string `gorm:"type:longtext"`
	Value      string `gorm:"type:longtext"`
	TaxonID    uint

	// Database references
	JobReportID uint `gorm:"index"`
}

// JobResultAccessor hosts every job results related dtabase method
type JobResultAccessor struct {
	*gorm.DB
}

// NewJobResultAccessor instanciates a new JobResultAccessor
func NewJobResultAccessor(db *gorm.DB) *JobResultAccessor {
	return &JobResultAccessor{db}
}

// Save helps to persists modifications on a job result
func (accessor *JobResultAccessor) Save(result *JobResult) error {
	return accessor.DB.Save(result).Error
}

// GetByID returns a single result from its ID
func (accessor *JobResultAccessor) GetByID(id uint) (*JobResult, error) {
	result := JobResult{}
	err := accessor.DB.First(&result, id).Error
	return &result, err
}

// GetAll returns every stored results
func (accessor *JobResultAccessor) GetAll() ([]*JobResult, error) {
	results := []*JobResult{}
	err := accessor.DB.Find(&results).Error
	return results, err
}

// GetByState returns every results with the given state
func (accessor *JobResultAccessor) GetByState(state string) ([]*JobResult, error) {
	results := []*JobResult{}
	err := accessor.DB.Where("state = ?", state).Find(&results).Error
	return results, err
}

// GetByStateAndReportID returns every results for the given report ID, filtered by state
func (accessor *JobResultAccessor) GetByStateAndReportID(reportID uint, state string) ([]*JobResult, error) {
	results := []*JobResult{}
	err := accessor.DB.Where("job_report_id = ? AND state = ?", reportID, state).Find(&results).Error
	return results, err
}
