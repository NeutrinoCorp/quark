package quark

import (
	"strconv"
	"time"

	"github.com/Shopify/sarama"
)

// MarshalKafkaMessage parses the given Message into a Apache Kafka producer message
func MarshalKafkaMessage(msg *Message) *sarama.ProducerMessage {
	partition, err := strconv.ParseInt(msg.Metadata.ExternalData[HeaderKafkaPartition], 10, 32)
	if err != nil {
		partition = 0
	}
	offset, err := strconv.ParseInt(msg.Metadata.ExternalData[HeaderKafkaOffset], 10, 64)
	if err != nil {
		offset = 0
	}

	return &sarama.ProducerMessage{
		Topic:     msg.Kind,
		Key:       nil,
		Value:     msg,
		Headers:   MarshalKafkaHeaders(msg),
		Metadata:  nil,
		Offset:    offset,
		Partition: int32(partition),
		Timestamp: time.Time{},
	}
}

// MarshalKafkaHeaders parses the given Message and its metadata into Apache Kafka's header types
func MarshalKafkaHeaders(msg *Message) []sarama.RecordHeader {
	h := make([]sarama.RecordHeader, 0)
	publishTime, err := msg.PublishTime.MarshalBinary()
	if err != nil {
		publishTime = []byte(msg.PublishTime.String())
	}

	h = append(h, sarama.RecordHeader{
		Key:   []byte(HeaderMessageId),
		Value: []byte(msg.Id),
	}, sarama.RecordHeader{
		Key:   []byte(HeaderMessageKind),
		Value: []byte(msg.Kind),
	}, sarama.RecordHeader{
		Key:   []byte(HeaderMessagePublishTime),
		Value: publishTime,
	}, sarama.RecordHeader{
		Key:   []byte(HeaderMessageCorrelationId),
		Value: []byte(msg.Metadata.CorrelationId),
	}, sarama.RecordHeader{
		Key:   []byte(HeaderMessageHost),
		Value: []byte(msg.Metadata.Host),
	}, sarama.RecordHeader{
		Key:   []byte(HeaderMessageRedeliveryCount),
		Value: []byte(strconv.Itoa(msg.Metadata.RedeliveryCount)),
	})
	for k, v := range msg.Metadata.ExternalData {
		h = append(h, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}
	return h
}

// UnmarshalKafkaHeaders parses the given Apache Kafka message into a Message
func UnmarshalKafkaHeaders(msgC *sarama.ConsumerMessage) *Message {
	msg := new(Message)
	msg.Metadata.ExternalData = map[string]string{}
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
				continue
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
	msg.Attributes = msgC.Value
	return msg
}
