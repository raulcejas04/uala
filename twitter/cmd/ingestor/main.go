package main

import (
	"twitter/pkg/routes"
)

const (
	TOPIC = "services"
)

func main() {
	routes.HandleRequests(TOPIC)
}
