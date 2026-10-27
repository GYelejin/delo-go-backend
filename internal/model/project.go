package model

import "time"

type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

// Project - сущность проекта по ТЗ ДЕЛО
type Project struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Owner       string    `json:"owner"`
	CreatedAt   time.Time `json:"created_at"`
}

// Task - задача внутри проекта/папки (из ДЗ 2 и 3)
type Task struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Title     string    `json:"title"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// SetStatus - метод с pointer receiver по требованию занятия 2
func (t *Task) SetStatus(newStatus Status) {
	t.Status = newStatus
}
