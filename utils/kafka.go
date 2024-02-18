package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func PostMessage(data interface{}, topic string, broker string, partition int32) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"" +
		"bootstrap.servers": broker,
		"acks": "all"})

	if err != nil {
		panic(err)
	}
	deliveryChan := make(chan kafka.Event, 10000)

	defer p.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          jsonData,
	}, deliveryChan)
	if err != nil {
		panic(err)
	}

	p.Flush(2 * 1000)
}

func SubscribeHandlerToTopic(topic string, partition int32, handler func([]byte), broker string, group string) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          group,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	partitions := []kafka.TopicPartition{{Topic: &topic, Partition: partition}}
	err = c.Assign(partitions)
	if err != nil {
		panic(err)
	}
	msgCount := 0
	err = c.Subscribe(topic, nil)
	if err != nil {
		panic(err)
	}

	for {
		ev := c.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			msgCount += 1
			if msgCount%3 == 0 {
				c.Commit()
			}
			handler(e.Value)

		case kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", e)
		case kafka.Error:
			_, err := fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			if err != nil {
				panic(err)
			}
			return
		}
	}
}
