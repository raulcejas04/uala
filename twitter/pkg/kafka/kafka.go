package kafka

import (
	"fmt"
	//"os"
	"sync"
	"time"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
)

var (
	//ctx context.Context
	lockKafka = &sync.Mutex{}
)

const (
	brokerAddress = "172.17.0.1:9094"
)

var Writer *kafka.Producer
var Reader *kafka.Consumer

func InitReader(topic string) {

	if Reader == nil {
		lockKafka.Lock()
		defer lockKafka.Unlock()

		log.Info("Inicializo reader ", brokerAddress)
		//l := log.New(os.Stdout, "kafka reader: ", 0)

		var err error

		Reader, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":  brokerAddress,
			"group.id":           "twitts",
			"group.instance.id":  "1",
			"enable.auto.commit": false,
			"auto.offset.reset":  "latest",
		})

		if err != nil {
			panic("could not open reader " + err.Error())
		} else {
			fmt.Println("Reader opened ")
		}

		Reader.SubscribeTopics([]string{topic}, nil)
		fmt.Println("Topic suscribed ", topic)
	}
}

func InitWriter(topic string) {

	if Writer == nil {
		lockKafka.Lock()
		defer lockKafka.Unlock()

		fmt.Println("Inicializo writer")
		//l := log.New(os.Stdout, "kafka writer: ", 0)
		//l := log.New()
		var err error

		Writer, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": brokerAddress,
			"acks":              "all",
		})

		if err != nil {
			panic("could not open producer " + err.Error())
		}

	}

}

func Close() {
	Reader.Close()
	Writer.Close()
}

func Produce(topic string, message string) {

	// intialize the writer with the broker addresses, and the topic
	InitWriter(topic)

	dt := time.Now()

	err := Writer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(dt.Format("20060201150405")),
		Value:          []byte(message),
	}, nil)

	if err != nil {
		panic("could not write message " + err.Error())
	}

	//log.Printf("Produce: %s", message )
}

func Consume(topic string) (string, string) {

	InitReader(topic)

	//fmt.Println( "entro reader" )
	msg, err := Reader.ReadMessage(-1)
	//fmt.Println( "paso reader" )
	if err != nil {
		panic("could not read message " + err.Error())
	}
	// after receiving the message, log its value
	fmt.Println("received key: ", string(msg.Key), " received message: ", string(msg.Value))

	return string(msg.Key), string(msg.Value)

}
