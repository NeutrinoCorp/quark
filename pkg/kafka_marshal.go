package pkg

import (
	"strconv"
	"time"

	"github.com/Shopify/sarama"
)

func UnmarshalKafkaHeaders(msgC *sarama.ConsumerMessage) *Message {
	msg := new(Message)
	for _, f := range msgC.Headers {
		switch string(f.Key) {
		case HeaderMessageId:
			msg.Id = string(f.Value)
		case HeaderMessageKind:
			msg.Kind = string(f.Value)
		case HeaderMessagePublishTime:
			t := time.Time{}
			if err := t.UnmarshalBinary(f.Value); err == nil {
				msg.PublishTime = t
			}
			msg.PublishTime = time.Time{} // set default empty value
		case HeaderMessageAttributes:
			msg.Attributes = f.Value
		case HeaderMessageCorrelationId:
			msg.Metadata.CorrelationId = string(f.Value)
		case HeaderMessageHost:
			msg.Metadata.Host = string(f.Value)
		case HeaderMessageRedeliveryCount:
			if r, err := strconv.Atoi(string(f.Value)); err == nil {
				msg.Metadata.RedeliveryCount = r
			}
		default:
			msg.Metadata.ExternalData[string(f.Key)] = string(f.Value)
		}
	}
	return msg
}
