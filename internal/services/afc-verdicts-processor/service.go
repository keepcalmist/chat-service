package afcverdictsprocessor

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/dgrijalva/jwt-go"
	"github.com/segmentio/kafka-go"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	clientmessageblockedjob "github.com/keepcalmist/chat-service/internal/services/outbox/jobs/client-message-blocked"
	clientmessagesentjob "github.com/keepcalmist/chat-service/internal/services/outbox/jobs/client-message-sent"
	"github.com/keepcalmist/chat-service/internal/types"
)

var errConst = errors.New("constant error")

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=afcverdictsprocessormocks

type messagesRepository interface {
	BlockMessage(ctx context.Context, msgID types.MessageID) error
	MarkAsVisibleForManager(ctx context.Context, msgID types.MessageID) error
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	backoffInitialInterval time.Duration `default:"100ms" validate:"min=50ms,max=1s"`
	backoffMaxElapsedTime  time.Duration `default:"5s" validate:"min=500ms,max=1m"`

	brokers         []string `option:"mandatory" validate:"min=1"`
	consumers       int      `option:"mandatory" validate:"min=1,max=16"`
	consumerGroup   string   `option:"mandatory" validate:"required"`
	verdictsTopic   string   `option:"mandatory" validate:"required"`
	verdictsSignKey string

	readerFactory KafkaReaderFactory `option:"mandatory" validate:"required"`
	dlqWriter     KafkaDLQWriter     `option:"mandatory" validate:"required"`

	txtor   transactor         `option:"mandatory" validate:"required"`
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
	outBox  outboxService      `option:"mandatory" validate:"required"`
}

type Service struct {
	Options
	logger *zap.Logger

	unmarshaler func(data []byte) (*verdict, error)
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	unmarshaler := unmarshalVerdictFromJSON

	if opts.verdictsSignKey != "" {
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(opts.verdictsSignKey))
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %w", err)
		}

		verdictUnmarshaler := func(data []byte) (*verdict, error) {
			return unmarshalVerdictWithKey(pubKey, data)
		}

		unmarshaler = verdictUnmarshaler
	}

	return &Service{
		Options:     opts,
		logger:      zap.L().Named("afcverdictsprocessor"),
		unmarshaler: unmarshaler,
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	defer s.dlqWriter.Close()

	eg := errgroup.Group{}

	for i := 0; i < s.consumers; i++ {
		eg.Go(func() error {
			return s.runConsumer(ctx)
		})
	}

	return eg.Wait()
}

func (s *Service) runConsumer(ctx context.Context) error {
	consumer := s.readerFactory(s.brokers, s.consumerGroup, s.verdictsTopic)
	defer consumer.Close()

	exponentialBackOff := backoff.NewExponentialBackOff()
	exponentialBackOff.InitialInterval = s.backoffInitialInterval
	exponentialBackOff.MaxElapsedTime = s.backoffMaxElapsedTime

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("context done, stopping afc_verdicts_processor")
			return nil
		default:
		}

		msg, err := consumer.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch message: %w", err)
		}

		var retryErr error

		retry := func() {
			ticker := backoff.NewTicker(exponentialBackOff)
			defer ticker.Stop()

			for range ticker.C {
				processingErr := s.processMessage(ctx, msg)
				if processingErr != nil {
					retryErr = multierr.Append(retryErr, processingErr)
					if errors.Is(err, errConst) {
						s.logger.Error("failed to processMessage message", zap.Error(processingErr))
						break
					}

					s.logger.Error("failed to processMessage message", zap.Error(processingErr))
					continue
				}

				break
			}
		}

		retry()

		if retryErr != nil {
			dlqMsg := msg

			dlqMsg.Headers = append(msg.Headers, kafka.Header{
				Key:   "LAST_ERROR",
				Value: []byte(retryErr.Error()),
			})
			dlqMsg.Headers = append(msg.Headers, kafka.Header{
				Key:   "ORIGINAL_PARTITION",
				Value: []byte(strconv.Itoa(msg.Partition)),
			})

			if err := s.dlqWriter.WriteMessages(ctx, dlqMsg); err != nil {
				return fmt.Errorf("failed to write to dlq: %w", err)
			}
		}

		if err = consumer.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("failed to commit message: %w", err)
		}
	}
}

func (s *Service) processMessage(ctx context.Context, msg kafka.Message) (err error) {
	gotVerdict, err := s.unmarshaler(msg.Value)
	if err != nil {
		return fmt.Errorf("%w unmarshal verdict: %s", errConst, err.Error())
	}

	switch gotVerdict.Status {
	case OK:
		if err := s.txtor.RunInTx(ctx, func(ctx context.Context) error {
			if err := s.msgRepo.MarkAsVisibleForManager(ctx, gotVerdict.MessageID); err != nil {
				return fmt.Errorf("mark as visible for manager err: %w", err)
			}

			if _, err := s.outBox.Put(ctx, clientmessagesentjob.Name, gotVerdict.MessageID.String(), time.Now()); err != nil {
				return fmt.Errorf("put job err: %w", err)
			}

			return nil
		}); err != nil {
			return fmt.Errorf("run in tx: %w", err)
		}
	case Suspicious:
		if err := s.txtor.RunInTx(ctx, func(ctx context.Context) error {
			if err := s.msgRepo.BlockMessage(ctx, gotVerdict.MessageID); err != nil {
				return fmt.Errorf("block message: %s", err)
			}

			if _, err := s.outBox.Put(ctx, clientmessageblockedjob.Name, gotVerdict.MessageID.String(), time.Now()); err != nil {
				return fmt.Errorf("put job: %w", err)
			}

			return nil
		}); err != nil {
			return fmt.Errorf("run in tx: %w", err)
		}
	default:
		return fmt.Errorf("%w unknown verdict status: %s", errConst, gotVerdict.Status)
	}

	return nil
}
