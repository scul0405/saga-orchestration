package eventhandler

import (
	"context"
	"encoding/json"
	"github.com/avast/retry-go"
	"github.com/scul0405/saga-orchestration/cmd/product/config"
	"github.com/scul0405/saga-orchestration/internal/common"
	"github.com/scul0405/saga-orchestration/internal/product/app"
	kafkaClient "github.com/scul0405/saga-orchestration/pkg/kafka"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/timeconvert"
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
	cfg        *config.Config
	logger     logger.Logger
	consumer   kafkaClient.ConsumerGroup
	producer   kafkaClient.Producer
	productSvc app.Application
}

func NewEventHandler(
	cfg *config.Config,
	logger logger.Logger,
	consumer kafkaClient.ConsumerGroup,
	producer kafkaClient.Producer,
	productSvc app.Application) EventHandler {
	return &eventHandler{
		cfg:        cfg,
		logger:     logger,
		consumer:   consumer,
		producer:   producer,
		productSvc: productSvc,
	}
}

func (h *eventHandler) Run(ctx context.Context) {
	go h.consumer.ConsumeTopic(ctx, poolSize, common.UpdateProductInventoryGroupID, common.UpdateProductInventoryTopic, h.updateProductInventoryWorker)
	go h.consumer.ConsumeTopic(ctx, poolSize, common.RollbackProductInventoryGroupID, common.RollbackProductInventoryTopic, h.rollbackProductInventoryWorker)
}

func (h *eventHandler) updateProductInventoryWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			h.logger.Errorf("Product.UpdateProductInventoryWorker: FetchMessage", err)
			return
		}

		var purchase pb.CreatePurchaseRequest
		if err = json.Unmarshal(m.Value, &purchase); err != nil {
			h.logger.Errorf("Product.UpdateProductInventoryWorker: UnmarshalProto", err)
			continue
		}
		h.logger.Infof("UpdateProductInventoryWorker: %v, message at topic/partition/offset %v/%v/%v: %s = %s\n", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		reply := pb.CreatePurchaseResponse{
			PurchaseId: purchase.PurchaseId,
		}
		if err = retry.Do(func() error {
			err = h.productSvc.Commands.UpdateProductInventory.Handle(ctx, decodePb2UpdateProductInventoryCmd(&purchase))
			if err != nil {
				reply.Success = false
				reply.ErrorMessage = err.Error()
			} else {
				reply.Success = true
			}

			reply.Purchase = purchase.Purchase
			reply.Timestamp = timeconvert.Time2pbTimestamp(time.Now())
			payload, _ := json.Marshal(&reply)

			return h.producer.PublishMessage(ctx, kafka.Message{
				Topic: common.ReplyTopic,
				Value: payload,
				Headers: []kafka.Header{
					{
						Key:   common.HandlerHeader,
						Value: []byte(common.UpdateProductInventoryHandler),
					},
				},
			})
		},
			retry.Attempts(retryAttempts),
			retry.Delay(retryDelay),
			retry.Context(ctx),
		); err != nil {
			// TODO: publish error message to error topic
		}

		err = r.CommitMessages(ctx, m)
		if err != nil {
			h.logger.Errorf("Product.UpdateProductInventoryWorker: CommitMessages", err)
		}
	}
}

func (h *eventHandler) rollbackProductInventoryWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			h.logger.Errorf("Product.RollbackProductInventoryWorker: FetchMessage", err)
			return
		}

		var purchase pb.CreatePurchaseRequest
		if err = json.Unmarshal(m.Value, &purchase); err != nil {
			h.logger.Errorf("Product.RollbackProductInventoryWorker: UnmarshalProto", err)
			continue
		}
		h.logger.Infof("RollbackProductInventoryWorker: %v, message at topic/partition/offset %v/%v/%v: %s = %s\n", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		reply := pb.CreatePurchaseResponse{
			PurchaseId: purchase.PurchaseId,
		}
		if err = retry.Do(func() error {
			err = h.productSvc.Commands.RollbackProductInventory.Handle(ctx, decodePb2RollbackProductInventoryCmd(&purchase))
			if err != nil {
				reply.Success = false
				reply.ErrorMessage = err.Error()
			} else {
				reply.Success = true
			}

			reply.Purchase = purchase.Purchase
			reply.Timestamp = timeconvert.Time2pbTimestamp(time.Now())
			payload, _ := json.Marshal(&reply)

			return h.producer.PublishMessage(ctx, kafka.Message{
				Topic: common.ReplyTopic,
				Value: payload,
				Headers: []kafka.Header{
					{
						Key:   common.HandlerHeader,
						Value: []byte(common.RollbackProductInventoryHandler),
					},
				},
			})
		},
			retry.Attempts(retryAttempts),
			retry.Delay(retryDelay),
			retry.Context(ctx),
		); err != nil {
			// TODO: publish error message to error topic
		}

		err = r.CommitMessages(ctx, m)
		if err != nil {
			h.logger.Errorf("Product.RollbackProductInventoryWorker: CommitMessages", err)
		}
	}
}
