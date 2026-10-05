package storage

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"time"

	"hw1taskmanager/internal/domain"
	"hw1taskmanager/internal/task"
	"hw1taskmanager/internal/validation"
)

type Validator interface {
	ValidateTaskData(desc string, deadline time.Time) error
	CanChangeStatus(status task.Status) error
}

type Storage struct {
	validator  Validator
	nextID     task.TaskID
	savingPath string
	tasks      map[task.TaskID]*task.Task
}

func NewStorage() *Storage {
	s := &Storage{
		validator: validation.NewValidator(),
		tasks: make(map[task.TaskID]*task.Task),
	}
	s.Restore()
	return s
}

func (s *Storage) AddTask(desc string, deadline time.Time) (*task.Task, error) {
	err := s.validator.ValidateTaskData(desc, deadline)
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	newTask := task.NewTask(s.nextID, desc, deadline)
	s.tasks[s.nextID] = newTask
	s.nextID++

	err = s.SaveTasks()
	if err != nil {
		delete(s.tasks, s.nextID)
		s.nextID--
		return nil, fmt.Errorf("failed to save tasks: %w", err)
	}
	return newTask, nil
}

func (s *Storage) GetTasks() ([]*task.Task, error) {
	tasks := slices.Collect(maps.Values(s.tasks))
	sort.Slice(tasks, func(i, j int) bool {
		if !tasks[i].Deadline.Equal(tasks[j].Deadline) {
			return tasks[i].Deadline.Before(tasks[j].Deadline)
		}

		return tasks[i].ID < tasks[j].ID
	})

	return tasks, nil
}

func (s *Storage) GetTaskByID(id task.TaskID) (*task.Task, error) {
	targetTask, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("ID %d: %w", id, domain.ErrTaskNotFound)
	}
	return targetTask, nil
}

func (s *Storage) ChangeTaskStatus(id task.TaskID, newStatus task.Status) (*task.Task, error) {
	targetTask, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("ID %d: %w", id, domain.ErrTaskNotFound)
	}

	err := s.validator.CanChangeStatus(targetTask.Status)
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	oldStatus := targetTask.Status
	targetTask.Status = newStatus

	err = s.SaveTasks()
	if err != nil {
		targetTask.Status = oldStatus
		return nil, fmt.Errorf("failed to save tasks: %w", err)
	}

	return targetTask, nil
}

func (s *Storage) ChangeTaskData(id task.TaskID, field string, newData task.TaskDTO) (*task.Task, error) {
	targetTask, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("ID %d: %w", id, domain.ErrTaskNotFound)
	}
	
	switch field {
		case "description":
			newData.Deadline = targetTask.Deadline
		case "deadline":
			newData.Description = targetTask.Description
	}
	
	err := s.validator.ValidateTaskData(newData.Description, newData.Deadline)
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	oldTask := targetTask
	targetTask.Description = newData.Description
	targetTask.Deadline = newData.Deadline

	err = s.SaveTasks()
	if err != nil {
		targetTask = oldTask
		return nil, fmt.Errorf("failed to save tasks: %w", err)
	}

	return targetTask, nil
}

func (s *Storage) DeleteTask(id task.TaskID) error {
	_, ok := s.tasks[id]
	if !ok {
		return fmt.Errorf("ID %d: %w", id, domain.ErrTaskNotFound)
	}

	delete(s.tasks, id)
	return nil
}
