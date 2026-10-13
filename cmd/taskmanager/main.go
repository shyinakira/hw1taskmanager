package main

import (
	"hw1taskmanager/internal/client"
	"hw1taskmanager/internal/service"
	"hw1taskmanager/internal/storage"
	"hw1taskmanager/internal/validation"
	"log"
)

func main() {
	storage, err := storage.NewStorage()
	if err != nil {
		log.Fatalf("can't initialize storage: %v", err)
	}
	validator := validation.NewValidator()
	taskService := service.NewTasksService(storage, validator)
	c := client.NewClient(taskService)
	c.StartClient()
}
