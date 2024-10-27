package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
	configs "twitter/pkg/configs"
	kafka "twitter/pkg/kafka"
	redis "twitter/pkg/redis"

	log "github.com/sirupsen/logrus"
)

const (
	TOPIC = "services"
)

type producer struct {
	twitts chan string
	quit   chan chan error
}

func saveRedis(mensaje string, current time.Time, wg *sync.WaitGroup) {
	var mensajes []string
	log.Infof(" mensaje ", mensaje)
	partes := strings.Split(mensaje, ";")
	mensajes = append(mensajes, fmt.Sprintf("%s;%s", current.Format(time.RFC3339), partes[1]))
	redis.SAdd(partes[0], mensajes)
	wg.Done()
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
			current := time.Now().Local()
			go saveRedis(s, current, &wg)
		}
	}

}
