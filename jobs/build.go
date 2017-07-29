package jobs

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	StatusQueue   = "queue"
	StatusProceed = "proceed"
	StatusSuccess = "success"
	StatusFailure = "failure"
)

type Build struct {
	ID        int
	Number    int
	Status    string
	UserID    int
	JobID     int
	Log       string
	Params    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Build) TableName() string {
	return "job_build"
}

func (m Model) build(db *gorm.DB, params map[string]string, userID int) error {
	bts, err := json.Marshal(params)
	if err != nil {
		return err
	}

	b := Build{
		Params:    string(bts),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		JobID:     m.ID,
		Status:    StatusQueue,
		Number:    0,
	}

	rows, err := db.Raw(fmt.Sprintf("SELECT max(number) FROM %s WHERE job_id = %d", b.TableName(), m.ID)).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&b.Number)
	}
	b.Number++

	if err := db.Create(&b).Error; err != nil {
		return err
	}

	return nil
}

func ParseBuildQueue() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		//start build
	}
}
