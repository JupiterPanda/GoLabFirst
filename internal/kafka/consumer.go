package kafka

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/models"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

// The example could be found in franz-go consumers examples of DLQ realization

type Config struct {
	Brokers []string
	Topic   string
	GroupID string
}

func (c Config) validate() error {
	if len(c.Brokers) == 0 || c.Brokers[0] == "" {
		return errors.New("kafka: KAFKA_BROKERS is empty")
	}
	if c.Topic == "" {
		return errors.New("kafka: KAFKA_TOPIC is empty")
	}
	if c.GroupID == "" {
		return errors.New("kafka: KAFKA_GROUP_ID is empty (required for manual commits)")
	}
	return nil
}

// Consumer читает protobuf-события из Kafka и вызывает готовые методы UseCase.
type Consumer struct {
	client  *kgo.Client
	useCase UseCase
}

func NewConsumer(ctx context.Context, cfg Config, useCase UseCase) (*Consumer, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.GroupID),
		kgo.ConsumeTopics(cfg.Topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		kgo.DisableAutoCommit(),    // коммитим сами, после обработки
		kgo.BlockRebalanceOnPoll(), // ребаланс не влезает в середину батча
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka: create client: %w", err)
	}

	if err := client.Ping(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("kafka: ping cluster: %w", err)
	}

	return &Consumer{client: client, useCase: useCase}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	log.Println("Kafka consumer listening")
	defer log.Println("Kafka consumer stopped")

	for {
		fetches := c.client.PollRecords(ctx, 500)
		if fetches.IsClientClosed() || ctx.Err() != nil {
			return nil
		}

		fetches.EachError(func(topic string, partition int32, err error) {
			if !errors.Is(err, context.Canceled) {
				log.Printf("kafka: fetch error (topic=%s partition=%d): %v", topic, partition, err)
			}
		})

		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			if len(p.Records) == 0 {
				return
			}

			for _, rec := range p.Records {
				if err := c.handle(ctx, rec); err != nil {
					// сообщение пропускаем: офсет всё равно коммитится ниже,
					// иначе битое событие заблокировало бы всю партицию
					log.Printf("kafka: handling failed (partition=%d offset=%d): %v",
						rec.Partition, rec.Offset, err)
				}
			}

			if err := c.client.CommitRecords(ctx, p.Records...); err != nil {
				log.Printf("kafka: commit failed (partition=%d): %v", p.Partition, err)
			}
		})

		c.client.AllowRebalance() // обязательно при BlockRebalanceOnPoll
	}
}

func (c *Consumer) Close() { c.client.Close() }

type UseCase interface {
	RentBookByTitleAndReaderName(ctx context.Context, name, title string) error
	ReturnBookByTitleAndReaderName(ctx context.Context, name, title string) error
	CreateBook(ctx context.Context, book models.Book) error
	DeleteBook(ctx context.Context, id int) error
	AddCopyOfBookById(ctx context.Context, id int) error
	SubtractCopyOfBookById(ctx context.Context, id int) error
}
