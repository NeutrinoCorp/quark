package kafka

import (
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
)

// MarshalKafkaMessage parses the given Message into a Apache Kafka producer message
func MarshalKafkaMessage(msg *quark.Message) *sarama.ProducerMessage {
	partition, err := strconv.ParseInt(msg.Metadata.ExternalData[HeaderKafkaPartition], 10, 32)
	if err != nil {
		partition = 0
	}
	offset, err := strconv.ParseInt(msg.Metadata.ExternalData[HeaderKafkaOffset], 10, 64)
	if err != nil {
		offset = 0
	}

	return &sarama.ProducerMessage{
		Topic:     msg.Type,
		Key:       sarama.StringEncoder(msg.Id),
		Value:     msg,
		Headers:   MarshalKafkaHeaders(msg),
		Offset:    offset,
		Partition: int32(partition),
	}
}

// MarshalKafkaHeaders parses the given Message and its metadata into Apache Kafka's header types
func MarshalKafkaHeaders(msg *quark.Message) []sarama.RecordHeader {
	publishTime, err := msg.Time.MarshalBinary()
	if err != nil {
		publishTime = []byte(msg.Time.String())
	}

	h := make([]sarama.RecordHeader, 0)
	h = append(h, sarama.RecordHeader{
		Key:   []byte(quark.HeaderMessageId),
		Value: []byte(msg.Id),
	}, sarama.RecordHeader{
		Key:   []byte(quark.HeaderMessageType),
		Value: []byte(msg.Type),
	}, sarama.RecordHeader{
		Key:   []byte(quark.HeaderMessageSpecVersion),
		Value: []byte(msg.SpecVersion),
	}, sarama.RecordHeader{
		Key:   []byte(quark.HeaderMessageSource),
		Value: []byte(msg.Source),
	}, sarama.RecordHeader{
		Key:   []byte(quark.HeaderMessageDataContentType),
		Value: []byte(msg.ContentType),
	}, sarama.RecordHeader{
		Key:   []byte(quark.HeaderMessageDataSchema),
		Value: []byte(msg.DataSchema),
	}, sarama.RecordHeader{
		Key:   []byte(quark.HeaderMessageSubject),
		Value: []byte(msg.Subject),
	}, sarama.RecordHeader{
		Key:   []byte(quark.HeaderMessageTime),
		Value: publishTime,
	}, sarama.RecordHeader{
		Key:   []byte(quark.HeaderMessageCorrelationId),
		Value: []byte(msg.Metadata.CorrelationId),
	}, sarama.RecordHeader{
		Key:   []byte(quark.HeaderMessageHost),
		Value: []byte(msg.Metadata.Host),
	}, sarama.RecordHeader{
		Key:   []byte(quark.HeaderMessageRedeliveryCount),
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

// UnmarshalKafkaHeaders parses the given Apache Kafka headers into the given Quark Message
func UnmarshalKafkaHeaders(headers []*sarama.RecordHeader, msg *quark.Message) {
	msg.Metadata.ExternalData = map[string]string{}
	for _, f := range headers {
		switch string(f.Key) {
		case quark.HeaderMessageId:
			if msg.Id == "" {
				msg.Id = string(f.Value)
			}
		case quark.HeaderMessageType:
			msg.Type = string(f.Value)
		case quark.HeaderMessageSpecVersion:
			msg.SpecVersion = string(f.Value)
		case quark.HeaderMessageSource:
			msg.Source = string(f.Value)
		case quark.HeaderMessageDataContentType:
			msg.ContentType = string(f.Value)
		case quark.HeaderMessageDataSchema:
			msg.DataSchema = string(f.Value)
		case quark.HeaderMessageSubject:
			msg.Subject = string(f.Value)
		case quark.HeaderMessageTime:
			t := time.Time{}
			if err := t.UnmarshalBinary(f.Value); err == nil {
				msg.Time = t
			}
		case quark.HeaderMessageData:
			msg.Data = f.Value
		case quark.HeaderMessageCorrelationId:
			msg.Metadata.CorrelationId = string(f.Value)
		case quark.HeaderMessageHost:
			msg.Metadata.Host = string(f.Value)
		case quark.HeaderMessageRedeliveryCount:
			if r, err := strconv.Atoi(string(f.Value)); err == nil {
				msg.Metadata.RedeliveryCount = r
			}
		default:
			msg.Metadata.ExternalData[string(f.Key)] = string(f.Value)
		}
	}
}

// UnmarshalKafkaMessage parses the given Apache Kafka message into a Message
func UnmarshalKafkaMessage(msgKafka *sarama.ConsumerMessage, msg *quark.Message) {
	msg.Id = string(msgKafka.Key)
	msg.Data = msgKafka.Value
	UnmarshalKafkaHeaders(msgKafka.Headers, msg)
}
