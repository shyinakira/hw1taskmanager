package main

import "hw1taskmanager/internal/client"

func main() {
	c := client.NewClient()
	c.StartClient()
}
