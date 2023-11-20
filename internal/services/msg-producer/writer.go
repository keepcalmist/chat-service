package msgproducer

import (
	"math/big"

	"github.com/segmentio/kafka-go"

	"github.com/keepcalmist/chat-service/internal/logger"
)

const serviceName = "msg-producer"

func NewKafkaWriter(brokers []string, topic string, batchSize int) KafkaWriter {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     new(balancer),
		BatchSize:    batchSize,
		RequiredAcks: kafka.RequireOne,
		Async:        false,
		Logger:       logger.NewKafkaAdapted().WithServiceName(serviceName),
		ErrorLogger:  logger.NewKafkaAdapted().WithServiceName(serviceName).ForErrors(),
	}
}

type balancer struct{}

func (b *balancer) Balance(msg kafka.Message, partitions ...int) (partition int) {
	n := new(big.Int).SetBytes(msg.Key)
	result := new(big.Int).Mod(n, big.NewInt(int64(len(partitions))))
	return int(result.Int64())
}
