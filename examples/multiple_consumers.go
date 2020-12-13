package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/neutrinocorp/quark/internal"

	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
	"github.com/neutrinocorp/quark/brokers"
)

func ErrHandle(_ context.Context, topic, queue string, err error) {
	log.Print("ERR: At topic ", topic)
	log.Print("ERR: At queue ", queue)
	log.Print("ERR: ", err.Error())
}

func main() {
	broker := brokers.NewKafka([]string{"localhost:9092"}, newSaramaCfg())
	if err := broker.Register(*newConsumerNotifyUserOnMarketplaceSale(), *newConsumerNotifyOrgOnPayment()); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		log.Print("stopping")
		cancel()
	}()

	err := broker.ListenAndServe(ctx)
	if err != nil {
		log.Print("FATAL: ", err)
		cancel()
		return
	}
}

func newConsumerNotifyUserOnMarketplaceSale() *quark.Consumer {
	topicName := quark.FormatTopicName("neutrino", "marketplace", quark.MessageDomainEvent, "sale",
		"published", 1)
	queueName := quark.FormatQueueName("user", "notification", "notify_user", "sale_published")
	defConsumer := &internal.KafkaConsumerHandler{}

	return quark.NewConsumer(topicName, queueName, 3, defConsumer, ErrHandle,
		quark.NewRetry(topicName, queueName, 3, 1, defConsumer, ErrHandle),
		quark.NewDeadLetter(topicName, queueName, 1, defConsumer, ErrHandle))
}

func newConsumerNotifyOrgOnPayment() *quark.Consumer {
	topicName := quark.FormatTopicName("neutrino", "payment", quark.MessageDomainEvent, "product",
		"payed", 1)
	queueName := quark.FormatQueueName("organization", "notification", "notify_organization",
		"product_payed")
	defConsumer := &internal.KafkaConsumerHandler{}

	return quark.NewConsumer(topicName, queueName, 3, defConsumer, ErrHandle,
		quark.NewRetry(queueName, queueName, 3, 1, defConsumer, ErrHandle),
		quark.NewDeadLetter(queueName, queueName, 1, defConsumer, ErrHandle))
}

func newSaramaCfg() *sarama.Config {
	config := sarama.NewConfig()
	config.ClientID = "neutrino-sample"
	config.Consumer.Return.Errors = true
	/*
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
		config.Producer.RequiredAcks = sarama.WaitForAll*/
	return config
}
