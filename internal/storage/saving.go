package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"hw1taskmanager/internal/domain"
)

const SAVE_FILE_NAME = "saved_tasks.json"

func (s *Storage) resolveSavePath() (string, error) {
	if s.savingPath != "" {
		return s.savingPath, nil
	}

	defaultPath, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve save path: %w", err)
	}

	return defaultPath, nil
}

func (s *Storage) SaveTasks() error {
	tasksJSON, err := json.Marshal(s.tasks)

	if err != nil {
		return fmt.Errorf("failed to encode tasks data: %w", err)
	}

	saveFolder, err := s.resolveSavePath()
	if err != nil {
		return err
	}
	saveFullPath := path.Join(saveFolder, SAVE_FILE_NAME)

	file, err := os.Create(saveFullPath)
	if err != nil {
		return fmt.Errorf("failed to create save file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(tasksJSON)
	if err != nil {
		return fmt.Errorf("failed write data: %w", err)
	}

	return nil
}

func (s *Storage) Restore() error {
	saveFolder, err := s.resolveSavePath()
	if err != nil {
		return err
	}
	saveFullPath := path.Join(saveFolder, SAVE_FILE_NAME)

	_, exists := os.Stat(saveFullPath)
	if exists != nil {
		return nil
	}

	file, err := os.Open(saveFullPath)
	if err != nil {
		return fmt.Errorf("failed to read saved data: %w", err)
	}
	defer file.Close()

	binTasks, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read saved data: %w", err)
	}

	err = json.Unmarshal(binTasks, &s.tasks)
	if err != nil {
		return fmt.Errorf("failed to decode tasks data: %w", err)
	}
	nextID := domain.TaskID(0)
	for _, task := range s.tasks {
		nextID = max(nextID, task.ID+1)
	}
	s.nextID = domain.TaskID(nextID)

	return nil
}
