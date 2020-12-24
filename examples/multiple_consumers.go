package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	//nil := &pkg.KafkaConsumerHandler{}

	return quark.NewConsumer(topicName, queueName, 3, nil, ErrHandle,
		quark.NewRetry(topicName, queueName, 3, 1, nil, ErrHandle),
		quark.NewDeadLetter(topicName, queueName, 1, nil, ErrHandle))
}

func newConsumerNotifyOrgOnPayment() *quark.Consumer {
	topicName := quark.FormatTopicName("neutrino", "payment", quark.MessageDomainEvent, "product",
		"payed", 1)
	queueName := quark.FormatQueueName("organization", "notification", "notify_organization",
		"product_payed")
	// nil := &pkg.KafkaConsumerHandler{}

	return quark.NewConsumer(topicName, queueName, 3, nil, ErrHandle,
		quark.NewRetry(queueName, queueName, 3, 1, nil, ErrHandle),
		quark.NewDeadLetter(queueName, queueName, 1, nil, ErrHandle))
}
