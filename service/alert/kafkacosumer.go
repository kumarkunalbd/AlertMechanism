package alert

import "github.com/confluentinc/confluent-kafka-go/kafka"

type KafkaConsumer struct {
	brokerId      string
	cleindId      string
	cosumerGroup  string
	consumertopic string
}

// A reader will be intia
func NewKafaConsumer(brokerId string, cleindId string, consumerGroup string, consumerTopic string) *KafkaConsumer {
	return &KafkaConsumer{
		brokerId:      brokerId,
		cleindId:      cleindId,
		cosumerGroup:  consumerGroup,
		consumertopic: consumerTopic,
	}
}

func (kc *KafkaConsumer) GetKafkareader() *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kc.brokerId,
		"group.id":          kc.cosumerGroup,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	c.SubscribeTopics([]string{kc.consumertopic}, nil)
	return c
}
