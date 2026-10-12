package client

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"hw1taskmanager/internal/domain"
	"hw1taskmanager/internal/parse"
	"hw1taskmanager/internal/storage"
	"hw1taskmanager/internal/task"
)

const HELP_MSG = `
cli application for task management

COMMANDS:
    - add task: add "<task_description>" "<deadline>"
    - list all tasks: all
    - get task by id: task <task_id>
    - update task: change <task_id> <description/deadline/all> "<new_data>"...
    - change task status: status <task_id> <new_status>
    - delete task: delete <task_id>

deadline format: "02-01-2006 15:04"

`

type Storage interface {
	AddTask(desc string, deadline time.Time) (*task.Task, error)
	GetTasks() ([]*task.Task, error)
	GetTaskByID(id task.TaskID) (*task.Task, error)
	ChangeTaskStatus(id task.TaskID, newStatus task.Status) (*task.Task, error)
	ChangeTaskData(id task.TaskID, field string, newData task.TaskDTO) (*task.Task, error)
	DeleteTask(id task.TaskID) error
}

type Client struct {
	rw      *bufio.ReadWriter
	storage Storage
}

func NewClient() *Client {
	return &Client{
		rw: bufio.NewReadWriter(
			bufio.NewReader(os.Stdin),
			bufio.NewWriter(os.Stdout),
		),
		storage: storage.NewStorage(),
	}
}

func (c *Client) readLine() ([]string, error) {
	l, readErr := c.rw.ReadString('\n')
	if readErr != nil {
		return []string{}, fmt.Errorf("failed to read command: %w", readErr)
	}

	parts, readErr := parse.ParseCommandLine(strings.TrimSpace(l))
	if readErr != nil {
		return []string{}, fmt.Errorf("failed to read command: %w", readErr)
	}

	if len(parts) == 0 {
		return []string{}, fmt.Errorf("failed to read command: %w", readErr)
	}

	return parts, nil
}

func (c *Client) processCommand() (error, bool) {
	defer c.rw.Flush()

	parts, err := c.readLine()
	if err != nil {
		return err, false
	}

	cmd := parts[0]

	switch cmd {
	case "add":
		if len(parts) != 3 {
			return fmt.Errorf(`%w. format: add "new task" "30.01.2077 12:00"`, domain.ErrWrongArgsCount), false
		}

		desc := parts[1]
		deadline := parts[2]
		err := c.add(desc, deadline)
		if err != nil {
			return fmt.Errorf("processing error: %w", err), false
		}
	case "all":
		if len(parts) != 1 {
			return fmt.Errorf(`%w. format: all`, domain.ErrWrongArgsCount), false
		}

		err := c.all()
		if err != nil {
			return fmt.Errorf("processing error: %w", err), false
		}
	case "task":
		if len(parts) != 2 {
			return fmt.Errorf(`%w. format: task 1111`, domain.ErrWrongArgsCount), false
		}

		id, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf(`%w. id must be a number`, domain.ErrInvalidArgument), false
		}
		err = c.task(id)
		if err != nil {
			return fmt.Errorf("processing error: %w", err), false
		}
	case "status":
		if len(parts) != 3 {
			return fmt.Errorf(`%w. format: status 1111 done`, domain.ErrWrongArgsCount), false
		}

		id, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf(`%w. id must be a number`, domain.ErrInvalidArgument), false
		}

		status := task.StringToStatus(parts[2])
		if status == task.UNKNOWN {
			return fmt.Errorf(`%w. valid statuses: planned, in progress, done, canceled`, domain.ErrInvalidArgument), false
		}

		err = c.status(id, status)
		if err != nil {
			return fmt.Errorf("processing error: %w", err), false
		}
	case "change":
		if len(parts) < 4 {
			return fmt.Errorf(`%w. format: change 1111 description "another task"`, domain.ErrWrongArgsCount), false
		}

		id, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf(`%w. id must be a number`, domain.ErrInvalidArgument), false
		}

		field := parts[2]
		switch field {
		case "all":
			if len(parts) != 5 {
				return fmt.Errorf(`%w. format: change 1111 all "another task" "30.01.2077 12:00"`, domain.ErrWrongArgsCount), false
			}
			err = c.change(id, field, parts[3], parts[4])
			if err != nil {
				return fmt.Errorf("processing error: %w", err), false
			}
		case "description":
			err = c.change(id, field, parts[3], "")
			if err != nil {
				return fmt.Errorf("processing error: %w", err), false
			}
		case "deadline":
			err = c.change(id, field, "", parts[3])
			if err != nil {
				return fmt.Errorf("processing error: %w", err), false
			}
		default:
			return fmt.Errorf(`%w. valid fields: description, deadline, all`, domain.ErrInvalidArgument), false
		}
	case "delete":
		if len(parts) != 2 {
			return fmt.Errorf(`%w. format: delete 1111`, domain.ErrWrongArgsCount), false
		}

		id, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf(`%w. id must be a number`, domain.ErrInvalidArgument), false
		}

		err = c.delete(id)
		if err != nil {
			return fmt.Errorf("processing error: %w", err), false
		}
	case "help":
		_, err := c.rw.Write([]byte(HELP_MSG))
		if err != nil {
			return fmt.Errorf("processing error: %w", err), false
		}
	case "exit":
		return nil, true
	}

	return nil, false
}

func (c *Client) StartClient() {
	for {
		err, isEnd := c.processCommand()
		if err != nil {
			log.Println(err.Error())
		}

		if isEnd {
			return
		}
	}
}
