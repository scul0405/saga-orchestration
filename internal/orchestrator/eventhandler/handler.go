package eventhandler

import (
	"context"
	"encoding/json"
	"github.com/avast/retry-go"
	"github.com/scul0405/saga-orchestration/cmd/orchestrator/config"
	"github.com/scul0405/saga-orchestration/internal/common"
	"github.com/scul0405/saga-orchestration/internal/orchestrator/app"
	"github.com/scul0405/saga-orchestration/internal/orchestrator/domain/aggregate"
	"github.com/scul0405/saga-orchestration/internal/orchestrator/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/orchestrator/domain/valueobject"
	kafkaClient "github.com/scul0405/saga-orchestration/pkg/kafka"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	pb "github.com/scul0405/saga-orchestration/proto"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

var (
	retryAttempts uint = 10
	retryDelay         = 1 * time.Second
	poolSize           = 16
)

type EventHandler interface {
	Run(ctx context.Context)
}

type eventHandler struct {
	cfg      *config.Config
	logger   logger.Logger
	consumer kafkaClient.ConsumerGroup
	app      app.App
}

func NewEventHandler(cfg *config.Config, logger logger.Logger, consumer kafkaClient.ConsumerGroup, app app.App) EventHandler {
	return &eventHandler{
		cfg:      cfg,
		logger:   logger,
		consumer: consumer,
		app:      app,
	}
}

func (h *eventHandler) Run(ctx context.Context) {
	go h.consumer.ConsumeTopic(ctx, poolSize, common.PurchaseGroupID, common.PurchaseTopic, h.createPurchaseWorker)
	go h.consumer.ConsumeTopic(ctx, poolSize, common.ReplyGroupID, common.ReplyTopic, h.replyWorker)
}

func (h *eventHandler) createPurchaseWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			h.logger.Errorf("Orchestrator.CreatePurchaseWorker: FetchMessage", err)
			return
		}

		var purchase pb.CreatePurchaseRequest
		if err = json.Unmarshal(m.Value, &purchase); err != nil {
			h.logger.Errorf("Orchestrator.CreatePurchaseWorker: UnmarshalProto", err)
			continue
		}
		h.logger.Infof("Orchestrator.CreatePurchaseWorker: %v, message at topic/partition/offset %v/%v/%v: %s = %s\n", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		orderItems := make([]entity.OrderItem, len(purchase.Purchase.Order.OrderItems))
		for i, item := range purchase.Purchase.Order.OrderItems {
			orderItems[i] = entity.OrderItem{
				ID:       item.ProductId,
				Quantity: item.Quantity,
			}
		}

		domainPurchase := aggregate.Purchase{
			ID: purchase.PurchaseId,
			Order: &entity.Order{
				OrderItems: &orderItems,
				CustomerID: purchase.Purchase.Order.CustomerId,
			},
			Payment: &valueobject.Payment{
				Amount:       purchase.Purchase.Payment.Amount,
				CurrencyCode: purchase.Purchase.Payment.CurrencyCode,
			},
		}

		if err = retry.Do(func() error {
			err = h.app.StartTransaction(ctx, &domainPurchase)
			if err != nil {
				h.logger.Errorf("Orchestrator.CreatePurchaseWorker: StartTransaction", err)
			}

			return err
		},
			retry.Attempts(retryAttempts),
			retry.Delay(retryDelay),
			retry.Context(ctx),
		); err != nil {
			// TODO: publish error message to error topic
		}

		err = r.CommitMessages(ctx, m)
		if err != nil {
			h.logger.Errorf("Orchestrator.CreatePurchaseWorker: CommitMessages", err)
		}
	}
}

func (h *eventHandler) replyWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			h.logger.Errorf("Orchestrator.ReplyWorker: FetchMessage", err)
			return
		}

		h.logger.Infof("ReplyWorker: %v, message at topic/partition/offset %v/%v/%v: %s = %s\n", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		if err = retry.Do(func() error {
			err = h.app.HandleReply(ctx, &m)
			if err != nil {
				h.logger.Errorf("Orchestrator.ReplyWorker: StartTransaction", err)
			}

			return err
		},
			retry.Attempts(retryAttempts),
			retry.Delay(retryDelay),
			retry.Context(ctx),
		); err != nil {
			// TODO: publish error message to error topic
		}

		err = r.CommitMessages(ctx, m)
		if err != nil {
			h.logger.Errorf("Orchestrator.ReplyWorker: CommitMessages", err)
		}
	}
}
