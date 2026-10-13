package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"hw1taskmanager/internal/domain"
	"hw1taskmanager/internal/parse"
)

const HELP_MSG = `
cli application for task management

COMMANDS:
    - add task: add "<title>" "<deadline>"
    - list all tasks: all
    - get task by id: task <task_id>
    - update task: change <task_id> <title/deadline/all> "<new_data>"...
    - change task status: status <task_id> <new_status>
    - delete task: delete <task_id>

deadline format: "02-01-2006 15:04"

`

type TasksService interface {
	AddTask(desc string, deadlineStr string) (domain.Task, error)
	AllTasks() ([]domain.Task, error)
	GetTask(id domain.TaskID) (domain.Task, error)
	ChangeStatus(id domain.TaskID, newStatus domain.Status) (domain.Task, error)
	ChangeTask(id domain.TaskID, field string, newDesc string, newDeadlineStr string) (domain.Task, error)
	DeleteTask(id domain.TaskID) error
}

type Client struct {
	rw      *bufio.ReadWriter
	service TasksService
}

func NewClient(service TasksService) *Client {
	return &Client{
		rw: bufio.NewReadWriter(
			bufio.NewReader(os.Stdin),
			bufio.NewWriter(os.Stdout),
		),
		service: service,
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

	return parts, nil
}

func (c *Client) processCommand() (error, bool) {
	defer c.rw.Flush()

	parts, err := c.readLine()
	if err != nil {
		return err, false
	}

	if len(parts) == 0 {
		return nil, false
	}

	cmd := parts[0]

	switch cmd {
	case "add":
		return c.add(parts), false
	case "all":
		return c.all(parts), false
	case "task":
		return c.task(parts), false
	case "status":
		return c.status(parts), false
	case "change":
		return c.change(parts), false
	case "delete":
		return c.delete(parts), false
	case "help":
		_, err := c.rw.Write([]byte(HELP_MSG))
		if err != nil {
			return fmt.Errorf("processing error: %w", err), false
		}
	case "exit":
		return nil, true
	default:
		return fmt.Errorf("processing error: %w", domain.ErrUnknownCommand), false
	}

	return nil, false
}

func (c *Client) StartClient() {
	for {
		err, isExit := c.processCommand()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Println(err.Error())
		}

		if isExit {
			return
		}
	}
}
