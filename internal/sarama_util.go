package internal

import (
	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
)

func ParseSaramaConsumer(handler interface{}) (sarama.ConsumerGroupHandler, error) {
	h, ok := handler.(sarama.ConsumerGroupHandler)
	if !ok {
		return nil, quark.ErrConsumerHandlerNotValid
	}

	return h, nil
}
