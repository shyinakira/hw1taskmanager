package storage

import (
	"fmt"
	"maps"
	"slices"
	"time"

	"hw1taskmanager/internal/domain"
)

type Storage struct {
	nextID     domain.TaskID
	savingPath string
	tasks      map[domain.TaskID]*domain.Task
}

func NewStorage() *Storage {
	s := &Storage{
		tasks: make(map[domain.TaskID]*domain.Task),
	}
	s.Restore()
	return s
}

func (s *Storage) AddTask(desc string, deadline time.Time) (domain.Task, error) {
	newTask := domain.NewTask(s.nextID, desc, deadline)
	s.tasks[s.nextID] = newTask
	s.nextID++

	err := s.SaveTasks()
	if err != nil {
		delete(s.tasks, s.nextID)
		s.nextID--
		return domain.Task{}, fmt.Errorf("failed to save tasks: %w", err)
	}
	return *newTask, nil
}

func (s *Storage) GetTasks() ([]domain.Task, error) {
	tasks := slices.Collect(maps.Values(s.tasks))
	tasksVal := make([]domain.Task, len(tasks))
	for i, task := range tasks {
		tasksVal[i] = *task
	}

	return tasksVal, nil
}

func (s *Storage) GetTaskByID(id domain.TaskID) (domain.Task, error) {
	targetTask, ok := s.tasks[id]
	if !ok {
		return domain.Task{}, fmt.Errorf("ID %d: %w", id, domain.ErrTaskNotFound)
	}
	return *targetTask, nil
}

func (s *Storage) ChangeTaskStatus(id domain.TaskID, newStatus domain.Status) (domain.Task, error) {
	targetTask, ok := s.tasks[id]
	if !ok {
		return domain.Task{}, fmt.Errorf("ID %d: %w", id, domain.ErrTaskNotFound)
	}

	oldStatus := targetTask.Status
	targetTask.Status = newStatus

	err := s.SaveTasks()
	if err != nil {
		targetTask.Status = oldStatus
		return domain.Task{}, fmt.Errorf("failed to save tasks: %w", err)
	}

	return *targetTask, nil
}

func (s *Storage) ChangeTaskData(id domain.TaskID, field string, newData domain.TaskDTO) (domain.Task, error) {
	targetTask, ok := s.tasks[id]
	if !ok {
		return domain.Task{}, fmt.Errorf("ID %d: %w", id, domain.ErrTaskNotFound)
	}

	switch field {
	case "description":
		newData.Deadline = targetTask.Deadline
	case "deadline":
		newData.Description = targetTask.Description
	}

	oldTask := targetTask
	targetTask.Description = newData.Description
	targetTask.Deadline = newData.Deadline

	err := s.SaveTasks()
	if err != nil {
		targetTask = oldTask
		return domain.Task{}, fmt.Errorf("failed to save tasks: %w", err)
	}

	return *targetTask, nil
}

func (s *Storage) DeleteTask(id domain.TaskID) error {
	task, ok := s.tasks[id]
	if !ok {
		return fmt.Errorf("ID %d: %w", id, domain.ErrTaskNotFound)
	}

	delete(s.tasks, id)
	err := s.SaveTasks()
	if err != nil {
		s.tasks[id] = task
		return fmt.Errorf("failed to save tasks: %w", err)
	}

	return nil
}
