package myRedis

import (
	"fmt"
	"sync"
	"time"

	//redis "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

func SaveRedis(username string, mensaje string, current time.Time, wg *sync.WaitGroup) {
	var mensajes []string
	log.Infof("username ", username, " mensaje ", mensaje)
	mensajes = append(mensajes, fmt.Sprintf("%s;%s", current.Format(time.RFC3339), mensaje))
	SAdd(username, mensajes)
	wg.Done()
}
