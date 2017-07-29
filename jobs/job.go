package jobs

import (
	"github.com/bolsunovskyi/scheduler/plugins"
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

type History struct {
	ID        int
	Action    string
	UserID    int
	Log       string
	CreatedAt time.Time
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
