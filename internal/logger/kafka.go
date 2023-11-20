package logger

import (
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var _ kafka.Logger = (*KafkaAdapted)(nil)

type KafkaAdapted struct {
	Logger     *zap.Logger `option:"mandatory"`
	isErrorLvl bool
}

func (k *KafkaAdapted) Printf(s string, i ...interface{}) {
	if k.isErrorLvl {
		k.Logger.Error(fmt.Sprintf(s, i...))
		return
	}

	k.Logger.Info(fmt.Sprintf(s, i...))
}

func NewKafkaAdapted() *KafkaAdapted {
	return &KafkaAdapted{
		Logger: zap.L(),
	}
}

func (k *KafkaAdapted) WithServiceName(name string) *KafkaAdapted {
	k.Logger = zap.L().Named(name)
	return k
}

func (k *KafkaAdapted) ForErrors() *KafkaAdapted {
	k.isErrorLvl = true
	return k
}
