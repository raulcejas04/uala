package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
	database "twitter/models/database"
	configs "twitter/pkg/configs"
	kafka "twitter/pkg/kafka"
	redis "twitter/pkg/redis"
	service "twitter/service"

	log "github.com/sirupsen/logrus"
)

const (
	TOPIC = "services"
)

type producer struct {
	twitts chan string
	quit   chan chan error
}

func (p *producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	res := <-ch
	return res
}

func main() {

	var wg sync.WaitGroup

	configs.InitConfig("./")
	redis.ConnectRedis()

	db := database.NewPostgresDB()
	err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	svc := service.NewDbService(db)

	prod := &producer{
		twitts: make(chan string),
		quit:   make(chan chan error),
	}

	fmt.Println("Recibe los mensajes nuevos y los pone en un canal para ser procesados")
	go func() {
		var k, m string

		var f bool = true

		for {

			if f {
				k, m = kafka.Consume(TOPIC)
			}

			fmt.Printf("mensaje leido %s %+v\n", k, m)

			if m == "Shutdown" {
				f = false
				k = ""
			}

			fmt.Println("Produce message: ", k, m)

			select {
			case prod.twitts <- m:
			case ch := <-prod.quit:
				close(prod.twitts)
				close(ch)
				return
			}
		}
	}()

	fmt.Println("Consume de un canal los mensajes nuevos y los procesa con distintos goroutines")
	for s := range prod.twitts {

		if s == "Shutdown" {
			fmt.Printf("Shutdown ")
			err := prod.Close()
			if err != nil {
				fmt.Printf("unexpected error: %v\n", err)
			}
		} else {

			wg.Add(1)

			partes := strings.Split(s, ";")
			result, err := svc.GetSeguidores(partes[0])
			if err != nil {
				log.Fatalf("Failed to fetch data: %v", err)
			}
			for _, seguidor := range result {
				current := time.Now().Local()
				go redis.SaveRedis(seguidor["username"].(string), partes[1], current, &wg)
			}
		}
	}

}
