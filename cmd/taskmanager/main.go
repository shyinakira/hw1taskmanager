package main

import (
	"hw1taskmanager/internal/client"
	"hw1taskmanager/internal/service"
	"hw1taskmanager/internal/storage"
)

func main() {
	storage := storage.NewStorage()
	taskService := service.NewTasksService(storage)
	c := client.NewClient(taskService)
	c.StartClient()
}
