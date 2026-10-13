package service

import (
	"sort"
	"time"

	"hw1taskmanager/internal/domain"
)

type Storage interface {
	AddTask(desc string, deadline time.Time) (domain.Task, error)
	GetTasks() ([]domain.Task, error)
	GetTaskByID(id domain.TaskID) (domain.Task, error)
	ChangeTaskStatus(id domain.TaskID, newStatus domain.Status) (domain.Task, error)
	ChangeTaskData(id domain.TaskID, field string, newData domain.TaskDTO) (domain.Task, error)
	DeleteTask(id domain.TaskID) error
}

type TasksService struct {
	storage Storage
}

func NewTasksService(s Storage) *TasksService {
	return &TasksService{
		storage: s,
	}
}

func convertTime(stringTime string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Yekaterinburg")
	if err != nil {
		return time.Time{}, domain.ErrTimezone
	}
	time, err := time.ParseInLocation(domain.TIME_LAYOUT, stringTime, loc)
	if err != nil {
		return time, err
	}

	return time, nil
}

func (s *TasksService) AddTask(desc string, deadlineStr string) (domain.Task, error) {
	deadline, err := convertTime(deadlineStr)
	if err != nil {
		return domain.Task{}, err
	}

	task, err := s.storage.AddTask(desc, deadline)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (s *TasksService) AllTasks() ([]domain.Task, error) {
	tasks, err := s.storage.GetTasks()
	if err != nil {
		return []domain.Task{}, err
	}

	sort.Slice(tasks, func(i, j int) bool {
		if !tasks[i].Deadline.Equal(tasks[j].Deadline) {
			return tasks[i].Deadline.Before(tasks[j].Deadline)
		}
		return tasks[i].ID < tasks[j].ID
	})
	return tasks, nil
}

func (s *TasksService) GetTask(id int) (domain.Task, error) {
	task, err := s.storage.GetTaskByID(domain.TaskID(id))
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (s *TasksService) ChangeStatus(id int, newStatus domain.Status) (domain.Task, error) {
	task, err := s.storage.ChangeTaskStatus(domain.TaskID(id), newStatus)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (s *TasksService) ChangeTask(id int, field string, newDesc string, newDeadlineStr string) (domain.Task, error) {
	var deadline time.Time
	var err error
	if newDeadlineStr != "" {
		deadline, err = convertTime(newDeadlineStr)
		if err != nil {
			return domain.Task{}, err
		}
	}

	task, err := s.storage.ChangeTaskData(domain.TaskID(id), field, domain.TaskDTO{
		Description: newDesc,
		Deadline:    deadline,
	})
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (s *TasksService) DeleteTask(id int) error {
	err := s.storage.DeleteTask(domain.TaskID(id))
	if err != nil {
		return err
	}

	return nil
}
