package jobs

import (
	"encoding/json"
	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type Model struct {
	ID                 int                 `json:"id"`
	Name               string              `json:"name" validate:"required"`
	Description        string              `json:"description"`
	IsEnabled          bool                `json:"is_enabled"`
	Params             []map[string]string `json:"params" gorm:"-"`
	PramsEncoded       string              `gorm:"column:params" json:"-"`
	RemoteTriggerToken string              `json:"remote_token"`
	PeriodicalSchedule string              `json:"schedule"`
	BuildPath          string              `json:"build_path"`
	StepsEncoded       string              `gorm:"column:steps" json:"-"`
	Steps              []Step              `gorm:"-" json:"steps"`
	CreatedAt          time.Time           `json:"-"`
	UpdatedAt          time.Time           `json:"-"`
	DeletedAt          *time.Time          `json:"-"`
}

func (Model) TableName() string {
	return "job"
}

type Step struct {
	PluginName  string              `json:"name"`
	Description string              `json:"description"`
	Params      []plugins.ItemParam `json:"schema"`
}

type Tab struct {
	ID   int
	Name string
}

type TabJon struct {
	TabID    int
	JobID    int
	Position int
}

func getJobByID(db *gorm.DB, idStr string) (*Model, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	var j Model
	if err := db.Where("id = ?", id).First(&j).Error; err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(j.StepsEncoded), &j.Steps)
	err = json.Unmarshal([]byte(j.PramsEncoded), &j.Params)
	if err != nil {
		return nil, err
	}

	return &j, nil
}
