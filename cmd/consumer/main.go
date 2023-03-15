package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gointensivo2/internal/infra/database"
	"github.com/gointensivo2/internal/usecase"
	"github.com/gointensivo2/pkg/kafka"
	"github.com/gointensivo2/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"

	// Sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, error := sql.Open("sqlite3", "./orders.db")
	if error != nil {
		panic(error)
	}
	defer db.Close()

	repository := database.NewOrderRepository(db)
	usecase := usecase.CalculateFinalPrice{OrderRepository: repository}

	msgChanKafka := make(chan *ckafka.Message)
	topics := []string{"orders"}
	servers := "host.docker.internal:9092"
	go kafka.Consume(topics, servers, msgChanKafka)
	go kafkaWorker(msgChanKafka, usecase)

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	msgChanRabbit := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, msgChanRabbit)
	rabbitmqWorker(msgChanRabbit, usecase)

}

func kafkaWorker(msgChan chan *ckafka.Message, uc usecase.CalculateFinalPrice) {
	for msg := range msgChan {
		var OrderInputDTO usecase.OrderInputDTO
		error := json.Unmarshal(msg.Value, &OrderInputDTO)
		if error != nil {
			panic(error)
		}
		outputDto, error := uc.Execute(OrderInputDTO)
		if error != nil {
			panic(error)
		}
		fmt.Println("Kafka Worker has processed order: %s\n", outputDto.ID)
	}
}

func rabbitmqWorker(msgChan chan amqp.Delivery, uc usecase.CalculateFinalPrice) {
	for msg := range msgChan {
		var OrderInputDTO usecase.OrderInputDTO
		error := json.Unmarshal(msg.Body, &OrderInputDTO)
		if error != nil {
			panic(error)
		}
		outputDto, error := uc.Execute(OrderInputDTO)
		if error != nil {
			panic(error)
		}
		msg.Ack(false)
		fmt.Println("RabbitMQ Worker has processed order: %s\n", outputDto.ID)
	}
}
