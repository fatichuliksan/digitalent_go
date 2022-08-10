package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Description string
	Assignee    string
	DeadlineAt  time.Time
	Status      string `gorm:"default:Pending"`
}
