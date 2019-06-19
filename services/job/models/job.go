package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	// JobStatusNew represents a new job that it has not been executed yet
	JobStatusNew = "NEW"
	// JobStatusIdle represents a finished job
	JobStatusIdle = "IDLE"
	// JobStatusPending represents a running job
	JobStatusPending = "PENDING"
	// JobStatusStacked represents a stacked job
	JobStatusStacked = "STACKED"
	// JobStatusArchived represents an archived job for quality and history purposes
	JobStatusArchived = "ARCHIVED"

	// JobTypeRedList indicates that a job is a red list job
	JobTypeRedList = "RL"
	// JobTypeZnieff indicates that a job is a ZNIEFF job
	JobTypeZnieff = "ZN"
	// JobTypeStatus indicates that a job is a regulatory status job
	JobTypeStatus = "ST"
)

// Job is the data model used to represent a job
type Job struct {
	ID         uint   `gorm:"primary_key"`
	Name       string `gorm:"type:varchar(2048)"`
	File       string `gorm:"type:varchar(2048)"`
	LastUpdate time.Time
	Status     string
	Type       string
	Data       string `gorm:"type:longtext"`
	Reports    []*JobReport
}

// JobAccessor hosts every job database methods
type JobAccessor struct {
	*gorm.DB
}

// NewJobAccessor allow to instanciate a new JobAccessor with the given database connextion
func NewJobAccessor(db *gorm.DB) *JobAccessor {
	return &JobAccessor{db}
}

// BeforeCreate helps to the the default status to a new job
func (job *Job) BeforeCreate() (err error) {
	job.Status = JobStatusNew
	return nil
}

// Create helps to persist a new job
func (accessor *JobAccessor) Create(job *Job) (*Job, error) {
	err := accessor.DB.Create(job).Error
	return job, err
}

// Save helps to save a job after it has been modified
func (accessor *JobAccessor) Save(job *Job) error {
	return accessor.DB.Save(job).Error
}

// GetAll returns every saved jobs
func (accessor *JobAccessor) GetAll() ([]Job, error) {
	jobs := []Job{}
	err := accessor.DB.Preload("Reports").Find(&jobs).Error
	return jobs, err
}

// GetByStatus returns every job with the given status
func (accessor *JobAccessor) GetByStatus(status string) ([]Job, error) {
	jobs := []Job{}
	err := accessor.DB.Where("jobs.status LIKE ?", status).Find(&jobs).Error
	return jobs, err
}

// GetByID returns a single job from the given ID
func (accessor *JobAccessor) GetByID(id uint) (*Job, error) {
	job := Job{}
	err := accessor.DB.First(&job, id).Error
	return &job, err
}
