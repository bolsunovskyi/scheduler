package jobs

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/smhouse/pi/db"
	"log"
	"sync"
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
	if rows.Next() {
		rows.Scan(&b.Number)
	}
	rows.Close()
	b.Number++

	if err := db.Create(&b).Error; err != nil {
		return err
	}

	return nil
}

func parseBuildQueue(db *gorm.DB) {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		//start build
		var builds []Build
		if err := db.Where("status = 'queue'").Group("job_id").Find(&builds).Error; err != nil {
			fmt.Println(err)
			continue
		}

		var wg sync.WaitGroup

		for _, build := range builds {
			var b Build
			if err := db.Where("status = 'queue' AND job_id = ?", build.JobID).
				Order("created_at").Limit(1).First(&b).Error; err != nil {
				log.Println(err)
				continue
			}

			wg.Add(1)
			go b.run(db, &wg)
		}

		wg.Wait()
	}
}

func (b Build) run(db *gorm.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	var job Model
	if err := db.Where("id = ?", b.JobID).First(&job).Error; err != nil {
		log.Println(err)
		return
	}

	if err := json.Unmarshal([]byte(job.StepsEncoded), &job.Steps); err != nil {
		log.Println(err)
		return
	}

}
