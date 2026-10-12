package client

import (
	"time"

	"hw1taskmanager/internal/domain"
	"hw1taskmanager/internal/task"
)

func convertTime(stringTime string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Yekaterinburg")
	if err != nil {
		return time.Time{}, nil
	}
	time, err := time.ParseInLocation(domain.TIME_LAYOUT, stringTime, loc)
	if err != nil {
		return time, err
	}

	return time, nil
}

func (c *Client) add(desc string, deadlineStr string) error {
	deadline, err := convertTime(deadlineStr)
	if err != nil {
		return err
	}

	task, err := c.storage.AddTask(desc, deadline)
	if err != nil {
		return err
	}

	c.rw.WriteString("task added")
	c.rw.WriteString(task.String())
	return nil
}

func (c *Client) all() error {
	tasks, err := c.storage.GetTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		_, err := c.rw.WriteString(task.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) task(id int) error {
	task, err := c.storage.GetTaskByID(task.TaskID(id))
	if err != nil {
		return err
	}

	c.rw.WriteString(task.String())
	return nil
}

func (c *Client) status(id int, newStatus task.Status) error {
	task, err := c.storage.ChangeTaskStatus(task.TaskID(id), newStatus)
	if err != nil {
		return err
	}

	c.rw.WriteString("status changed")
	c.rw.WriteString(task.String())
	return nil
}

func (c *Client) change(id int, field string, newDesc string, newDeadlineStr string) error {
	var deadline time.Time
	var err error
	if newDeadlineStr != "" {
		deadline, err = convertTime(newDeadlineStr)
		if err != nil {
			return err
		}
	}

	task, err := c.storage.ChangeTaskData(task.TaskID(id), field, task.TaskDTO{
		Description: newDesc,
		Deadline:    deadline,
	})
	if err != nil {
		return err
	}

	c.rw.WriteString("task changed")
	c.rw.WriteString(task.String())
	return nil
}

func (c *Client) delete(id int) error {
	err := c.storage.DeleteTask(task.TaskID(id))
	if err != nil {
		return err
	}

	c.rw.WriteString("task deleted\n")
	return nil
}
