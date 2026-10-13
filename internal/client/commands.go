package client

import (
	"fmt"
	"strconv"

	"hw1taskmanager/internal/domain"
)

func (c *Client) add(parts []string) error {
	if len(parts) != 3 {
		return fmt.Errorf(`%w. format: add "new task" "30.01.2077 12:00"`, domain.ErrWrongArgsCount)
	}

	desc := parts[1]
	deadline := parts[2]
	task, err := c.service.AddTask(desc, deadline)
	if err != nil {
		return fmt.Errorf("processing error: %w", err)
	}

	_, err = c.rw.Write([]byte(task.String()))
	if err != nil {
		return fmt.Errorf("processing error: %w", err)
	}

	return nil
}

func (c *Client) all(parts []string) error {
	if len(parts) != 1 {
		return fmt.Errorf(`%w. format: all`, domain.ErrWrongArgsCount)
	}

	tasks, err := c.service.AllTasks()

	res := make([]byte, 0)
	for _, task := range tasks {
		res = append(res, []byte(task.String())...)
	}

	_, err = c.rw.Write(res)
	if err != nil {
		return fmt.Errorf("processing error: %w", err)
	}

	return nil
}

func (c *Client) task(parts []string) error {
	if len(parts) != 2 {
		return fmt.Errorf(`%w. format: task 1111`, domain.ErrWrongArgsCount)
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf(`%w. id must be a number`, domain.ErrInvalidArgument)
	}

	task, err := c.service.GetTask(id)

	_, err = c.rw.Write([]byte(task.String()))
	if err != nil {
		return fmt.Errorf("processing error: %w", err)
	}

	return nil
}

func (c *Client) status(parts []string) error {
	if len(parts) != 3 {
		return fmt.Errorf(`%w. format: status 1111 done`, domain.ErrWrongArgsCount)
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf(`%w. id must be a number`, domain.ErrInvalidArgument)
	}

	status := domain.StringToStatus(parts[2])
	if status == domain.UNKNOWN {
		return fmt.Errorf(`%w. valid statuses: planned, in progress, done, canceled`, domain.ErrInvalidArgument)
	}

	task, err := c.service.ChangeStatus(id, status)

	_, err = c.rw.Write([]byte(task.String()))
	if err != nil {
		return fmt.Errorf("processing error: %w", err)
	}

	return nil
}

func (c *Client) change(parts []string) error {
	if len(parts) < 4 {
		return fmt.Errorf(`%w. format: change 1111 description "another task"`, domain.ErrWrongArgsCount)
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf(`%w. id must be a number`, domain.ErrInvalidArgument)
	}

	field := parts[2]
	switch field {
	case "all":
		if len(parts) != 5 {
			return fmt.Errorf(`%w. format: change 1111 all "another task" "30.01.2077 12:00"`, domain.ErrWrongArgsCount)
		}
		task, err := c.service.ChangeTask(id, field, parts[3], parts[4])
		if err != nil {
			return fmt.Errorf("processing error: %w", err)
		}

		_, err = c.rw.Write([]byte(task.String()))
		if err != nil {
			return fmt.Errorf("processing error: %w", err)
		}
	case "description":
		task, err := c.service.ChangeTask(id, field, parts[3], "")
		if err != nil {
			return fmt.Errorf("processing error: %w", err)
		}

		_, err = c.rw.Write([]byte(task.String()))
		if err != nil {
			return fmt.Errorf("processing error: %w", err)
		}
	case "deadline":
		task, err := c.service.ChangeTask(id, field, "", parts[4])
		if err != nil {
			return fmt.Errorf("processing error: %w", err)
		}

		_, err = c.rw.Write([]byte(task.String()))
		if err != nil {
			return fmt.Errorf("processing error: %w", err)
		}
	default:
		return fmt.Errorf(`%w. valid fields: description, deadline, all`, domain.ErrInvalidArgument)
	}

	return nil
}

func (c *Client) delete(parts []string) error {
	if len(parts) != 2 {
		return fmt.Errorf(`%w. format: delete 1111`, domain.ErrWrongArgsCount)
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf(`%w. id must be a number`, domain.ErrInvalidArgument)
	}

	err = c.service.DeleteTask(id)
	if err != nil {
		return fmt.Errorf("processing error: %w", err)
	}

	_, err = c.rw.Write([]byte("task deleted\n"))
	if err != nil {
		return fmt.Errorf("processing error: %w", err)
	}

	return nil
}
