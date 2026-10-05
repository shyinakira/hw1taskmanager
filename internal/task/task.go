package task

import (
	"fmt"
	"hw1taskmanager/internal/domain"
	"time"
)

type Status int

type Field int

func StringToStatus(s string) Status {
	switch s {
	case "planned":
		return PLANNED
	case "in progress":
		return IN_PROGRESS
	case "canceled":
		return CANCELED
	case "done":
		return DONE
	default:
		return UNKNOWN
	}
}

func (s Status) String() string {
	switch s {
	case PLANNED:
		return "planned"
	case IN_PROGRESS:
		return "in progress"
	case CANCELED:
		return "canceled"
	case DONE:
		return "done"
	default:
		return "unknown"
	}
}

const (
	UNKNOWN Status = iota
	PLANNED
	IN_PROGRESS
	CANCELED
	DONE

	UndefinedID TaskID = -1
)

type TaskID int64

type Task struct {
	ID          TaskID    `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	Deadline    time.Time `json:"deadline"`
}

type TaskDTO struct {
	Description string
	Deadline    time.Time
}

func (t Task) String() string {
	return fmt.Sprintf(
		"\nID: %d\nTask description: %s\nStatus: %s\nDeadline: %v\n\n",
		t.ID,
		t.Description,
		t.Status.String(),
		t.Deadline.Format(domain.TIME_LAYOUT),
	)
}

func NewTask(id TaskID, description string, deadline time.Time) *Task {
	return &Task{
		ID:          id,
		Description: description,
		Status:      PLANNED,
		Deadline:    deadline,
	}
}
