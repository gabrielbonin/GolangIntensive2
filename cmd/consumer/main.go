package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gointensivo2/internal/infra/database"
	"github.com/gointensivo2/internal/usecase"
	"github.com/gointensivo2/pkg/kafka"

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
	kafkaWorker(msgChanKafka, usecase)

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