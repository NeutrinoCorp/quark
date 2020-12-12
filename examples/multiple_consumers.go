package main

import (
	"context"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
	"github.com/neutrinocorp/quark/brokers"
	"github.com/neutrinocorp/quark/subscribers"
)

func ErrHandle(_ context.Context, topic, queue string, err error) {
	log.Print("ERR: At topic ", topic)
	log.Print("ERR: At queue ", queue)
	log.Print("ERR: ", err.Error())
}

func main() {
	broker := brokers.NewKafka([]string{"localhost:9092"}, sarama.NewConfig(), ErrHandle)
	if err := broker.Register(*newConsumerMarketplace(), *newConsumerPayment()); err != nil {
		log.Fatal(err)
	}

	err := broker.ListenAndServe(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func newConsumerMarketplace() *quark.Consumer {
	topicName := quark.FormatTopicName("neutrino", "marketplace", quark.MessageDomainEvent, "sale",
		"published", 1)
	queueName := quark.FormatQueueName("user", "notification", "notify_user", "sale_published")
	defConsumer := &subscribers.KafkaConsumerHandler{}

	return quark.NewConsumer(topicName, queueName, 3, defConsumer,
		quark.NewRetry(topicName, queueName, 3, 1, defConsumer),
		quark.NewDeadLetter(topicName, queueName, 1, defConsumer))
}

func newConsumerPayment() *quark.Consumer {
	topicName := quark.FormatTopicName("neutrino", "payment", quark.MessageDomainEvent, "product",
		"payed", 1)
	queueName := quark.FormatQueueName("organization", "notification", "notify_organization",
		"product_payed")
	defConsumer := &subscribers.KafkaConsumerHandler{}

	return quark.NewConsumer(topicName, queueName, 10, defConsumer,
		quark.NewRetry(topicName, queueName, 3, 1, defConsumer),
		quark.NewDeadLetter(topicName, queueName, 1, defConsumer))
}

func newSaramaCfg() *sarama.Config {
	config := sarama.NewConfig()
	config.ClientID = "test"
	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.Retention = time.Minute * 1
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.Retry.Max = 3
	config.Consumer.Group.Session.Timeout = time.Second * 10
	config.Consumer.Retry.Backoff = time.Second * 2
	config.Consumer.Group.Rebalance.Retry.Max = 4
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Producer.Return.Errors = false
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1
	config.Producer.Return.Successes = false
	config.Producer.Retry.Backoff = time.Millisecond * 100
	config.Producer.Retry.Max = 3
	config.Producer.RequiredAcks = sarama.WaitForAll
	return config
}
