package service

import (
	"fmt"
	"sort"
	"time"

	"hw1taskmanager/internal/domain"
)

type Validator interface {
	ValidateTaskData(title string, deadline time.Time) error
	CanChangeStatus(oldStatus, newStatus domain.Status) error
}

type Storage interface {
	AddTask(title string, deadline time.Time) (domain.Task, error)
	GetTasks() ([]domain.Task, error)
	GetTaskByID(id domain.TaskID) (domain.Task, error)
	ChangeTaskStatus(id domain.TaskID, newStatus domain.Status) (domain.Task, error)
	ChangeTaskData(id domain.TaskID, field string, newData domain.TaskDTO) (domain.Task, error)
	DeleteTask(id domain.TaskID) error
}

type TasksService struct {
	validator Validator
	storage   Storage
}

func NewTasksService(s Storage, v Validator) *TasksService {
	return &TasksService{
		validator: v,
		storage:   s,
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

func (s *TasksService) AddTask(title string, deadlineStr string) (domain.Task, error) {
	deadline, err := convertTime(deadlineStr)
	if err != nil {
		return domain.Task{}, fmt.Errorf("error while trying to add task: %w", err)
	}

	err = s.validator.ValidateTaskData(title, deadline)
	if err != nil {
		return domain.Task{}, fmt.Errorf("error while trying to add task: %w", err)
	}

	task, err := s.storage.AddTask(title, deadline)
	if err != nil {
		return domain.Task{}, fmt.Errorf("error while trying to add task: %w", err)
	}

	return task, nil
}

func (s *TasksService) AllTasks() ([]domain.Task, error) {
	tasks, err := s.storage.GetTasks()
	if err != nil {
		return []domain.Task{}, fmt.Errorf("error while trying to list all task: %w", err)
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
		return domain.Task{}, fmt.Errorf("error while trying to get task by id: %w", err)
	}

	return task, nil
}

func (s *TasksService) ChangeStatus(id int, newStatus domain.Status) (domain.Task, error) {
	task, err := s.storage.GetTaskByID(domain.TaskID(id))
	if err != nil {
		return domain.Task{}, fmt.Errorf("error while trying to change task status: %w", err)
	}

	err = s.validator.CanChangeStatus(task.Status, newStatus)
	if err != nil {
		return domain.Task{}, fmt.Errorf("error while trying to change task status: %w", err)
	}

	task, err = s.storage.ChangeTaskStatus(domain.TaskID(id), newStatus)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (s *TasksService) ChangeTask(id int, field string, newTitle string, newDeadlineStr string) (domain.Task, error) {
	var deadline time.Time
	switch field {
	case "title":
		err := s.validator.ValidateTaskData(newTitle, deadline)
		if err != nil {
			return domain.Task{}, fmt.Errorf("error while trying to change task data: %w", err)
		}
	case "deadline":
		task, err := s.storage.GetTaskByID(domain.TaskID(id))
		if err != nil {
			return domain.Task{}, fmt.Errorf("error while trying to change task data: %w", err)
		}

		deadline, err = convertTime(newDeadlineStr)
		if err != nil {
			return domain.Task{}, fmt.Errorf("error while trying to change task data: %w", err)
		}

		err = s.validator.ValidateTaskData(task.Title, deadline)
		if err != nil {
			return domain.Task{}, fmt.Errorf("error while trying to change task data: %w", err)
		}
	case "all":
		var err error
		deadline, err = convertTime(newDeadlineStr)
		if err != nil {
			return domain.Task{}, fmt.Errorf("error while trying to change task data: %w", err)
		}

		err = s.validator.ValidateTaskData(newTitle, deadline)
		if err != nil {
			return domain.Task{}, fmt.Errorf("error while trying to change task data: %w", err)
		}
	}

	task, err := s.storage.ChangeTaskData(domain.TaskID(id), field, domain.TaskDTO{
		Title:    newTitle,
		Deadline: deadline,
	})
	if err != nil {
		return domain.Task{}, fmt.Errorf("error while trying to change task data: %w", err)
	}

	return task, nil
}

func (s *TasksService) DeleteTask(id int) error {
	err := s.storage.DeleteTask(domain.TaskID(id))
	if err != nil {
		return fmt.Errorf("error while trying to delete task: %w", err)
	}

	return nil
}
