package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goproject/internal/delivery/libgrpc"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Envelope — общий конверт сообщения.
// Producer кладёт: {"action":"rent_book","payload":{...}}
type Envelope struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type Config struct {
	Brokers  []string
	Topic    string
	GroupID  string
	DLQTopic string // опционально: куда складывать необрабатываемые сообщения
}

// Consumer читает события из Kafka и вызывает уже готовые методы UseCase.
// Ни один gRPC-метод не переопределяется — используется тот же интерфейс libgrpc.UseCase,
// который получает GRPCServer.
type Consumer struct {
	cl       *kgo.Client
	useCase  libgrpc.UseCase
	dlqTopic string
}

func NewConsumer(ctx context.Context, cfg Config, useCase libgrpc.UseCase) (*Consumer, error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.GroupID),
		kgo.ConsumeTopics(cfg.Topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()), // новая группа читает топик с начала
		kgo.DisableAutoCommit(),                           // коммитим вручную, только после успешной обработки
		kgo.BlockRebalanceOnPoll(),                        // ребаланс не влезает в середину обработки батча
		kgo.AllowAutoTopicCreation(),                      // удобно для DLQ в dev-окружении
	)
	if err != nil {
		return nil, fmt.Errorf("kafka: create client: %w", err)
	}

	if err := cl.Ping(ctx); err != nil { // сразу проверяем связность с кластером
		cl.Close()
		return nil, fmt.Errorf("kafka: ping cluster: %w", err)
	}

	return &Consumer{cl: cl, useCase: useCase, dlqTopic: cfg.DLQTopic}, nil
}

func (c *Consumer) Close() { c.cl.Close() }

// Run — блокирующий цикл. Запускать в горутине, как HTTP gateway в main.go.
func (c *Consumer) Run(ctx context.Context) error {
	log.Println("Kafka consumer listening")

	for {
		// PollRecords вместо PollFetches — при BlockRebalanceOnPoll ограничиваем
		// размер батча, чтобы обработка не держала ребаланс слишком долго.
		fetches := c.cl.PollRecords(ctx, 500)

		if fetches.IsClientClosed() {
			log.Println("Kafka consumer stopped: client closed")
			return nil
		}
		if err := ctx.Err(); err != nil {
			log.Println("Kafka consumer stopped: context done")
			return nil
		}

		// Ошибки фетча логируем, но не падаем: retriable franz-go разруливает сам
		fetches.EachError(func(topic string, partition int32, err error) {
			if !errors.Is(err, context.Canceled) {
				log.Printf("kafka: fetch error (topic=%s partition=%d): %v", topic, partition, err)
			}
		})

		// Обрабатываем партиции по очереди — порядок внутри партиции сохраняется
		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			for _, rec := range p.Records {
				if err := c.handle(ctx, rec); err != nil {
					log.Printf("kafka: handling failed (topic=%s partition=%d offset=%d): %v",
						rec.Topic, rec.Partition, rec.Offset, err)
					c.toDLQ(ctx, rec, err)
				}
			}

			// Коммитим офсет партии после её обработки — и для успешных,
			// и для отправленных в DLQ, чтобы одно битое сообщение не блокировало партицию.
			if len(p.Records) > 0 {
				if err := c.cl.CommitRecords(ctx, p.Records...); err != nil {
					log.Printf("kafka: commit failed (topic=%s partition=%d offset=%d): %v",
						p.Topic, p.Partition, p.Records[len(p.Records)-1].Offset+1, err)
				}
			}
		})

		c.cl.AllowRebalance() // обязательно при BlockRebalanceOnPoll
	}
}

func (c *Consumer) handle(ctx context.Context, rec *kgo.Record) error {
	var env Envelope
	if err := json.Unmarshal(rec.Value, &env); err != nil {
		return fmt.Errorf("invalid envelope: %w", err)
	}

	handler, ok := c.routes()[env.Action]
	if !ok {
		return fmt.Errorf("unknown action %q", env.Action)
	}

	return handler(ctx, env.Payload)
}

// toDLQ пишет необработанное сообщение в отдельный топик тем же клиентом
func (c *Consumer) toDLQ(ctx context.Context, rec *kgo.Record, cause error) {
	if c.dlqTopic == "" {
		return
	}

	dead := &kgo.Record{
		Topic: c.dlqTopic,
		Key:   rec.Key,
		Value: rec.Value,
		Headers: []kgo.RecordHeader{
			{Key: "error", Value: []byte(cause.Error())},
			{Key: "origin-topic", Value: []byte(rec.Topic)},
			{Key: "origin-offset", Value: []byte(fmt.Sprint(rec.Offset))},
		},
	}

	if err := c.cl.ProduceSync(ctx, dead).FirstErr(); err != nil {
		log.Printf("kafka: failed to write to DLQ: %v", err)
	}
}
