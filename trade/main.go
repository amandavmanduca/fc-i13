package main

import (
	"fmt"
	"github.com/amandavmanduca/fc-i13/go/service/internal/market/transformer"
	"encoding/json"
	"github.com/amandavmanduca/fc-i13/go/service/internal/market/dto"
	"sync"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/amandavmanduca/fc-i13/go/service/internal/infra/kafka"
	"github.com/amandavmanduca/fc-i13/go/service/internal/market/entity"
)

func main() {
	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	kafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id": "myGroup",
		"auto.offset.reset": "earliest",
	}
	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"input"})

	go kafka.Consume(kafkaMsgChan) // T2

	book := entity.NewBook(ordersIn, ordersOut, wg)
	go book.Trade() //T3

	go func() {
		for msg := range kafkaMsgChan {
			wg.Add(1)
			fmt.Println(string(msg.Value))
			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput) // json -> objeto
			if err != nil {
				panic(err)
			}
			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", " ") // objeto -> json
		if err != nil {
			fmt.Println(err)
		}
		producer.Publish(outputJson, []byte("orders"), "output")
	}
}