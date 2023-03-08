package kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

func Consume(topics []string, servers string, msgChan chan *ckafka.Message) {
	kafkaConsumer, error := ckafka.NewConsumer(&ckafka.ConfigMap{
		"bootstrap.servers": servers,
		"group.id":          "goapp",
		"auto.offset.reset": "earliest",
	})

	if error != nil {
		panic(error)
	}

	error = kafkaConsumer.SubscribeTopics(topics, nil)
	if error != nil {
		panic(error)
	}

	for {
		msg, error := kafkaConsumer.ReadMessage(-1)
		if error == nil {
			msgChan <- msg
		}
	}
}
