package jobs

import "time"

type Model struct {
	ID                 int
	Name               string
	Description        string
	IsEnabled          bool
	IsParametrized     bool
	Prams              string
	RemoteTriggerToken string
	PeriodicalSchedule string
	BuildPath          string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          time.Time
	//DiscardOldBuilds     bool
	//DaysToKeepBuilds     int
	//MaxNOfBuilds         int
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
