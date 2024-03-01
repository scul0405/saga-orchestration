package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/scul0405/saga-orchestration/internal/common"
	"github.com/scul0405/saga-orchestration/internal/orchestrator/domain/aggregate"
	"github.com/scul0405/saga-orchestration/internal/orchestrator/domain/event"
	kafkaClient "github.com/scul0405/saga-orchestration/pkg/kafka"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/timeconvert"
	pb "github.com/scul0405/saga-orchestration/proto"
	"github.com/segmentio/kafka-go"
	"time"
)

type App interface {
	StartTransaction(ctx context.Context, purchase *aggregate.Purchase) error
	HandleReply(ctx context.Context, msg *kafka.Message) error
}

type app struct {
	logger   logger.Logger
	producer kafkaClient.Producer
}

func NewApp(logger logger.Logger, producer kafkaClient.Producer) App {
	return &app{
		logger:   logger,
		producer: producer,
	}
}

func (a *app) StartTransaction(ctx context.Context, purchase *aggregate.Purchase) error {
	err := a.publishPurchaseResult(ctx, &event.PurchaseResult{
		PurchaseID: purchase.ID,
		Status:     event.StatusExecute,
		Step:       event.StepUpdateProductInventory,
	})
	if err != nil {
		return err
	}

	// Call the product service to update the product inventory
	pbPurchase := encodeModel2PurchaseRequest(purchase)

	payload, _ := json.Marshal(pbPurchase)
	return a.producer.PublishMessage(ctx, kafka.Message{
		Topic: common.UpdateProductInventoryTopic,
		Value: payload,
	})
}

func (a *app) HandleReply(ctx context.Context, msg *kafka.Message) error {
	switch string(msg.Headers[0].Value) {
	case common.UpdateProductInventoryHandler:
		purchaseResult, err := decodePbResponseToEventModel(msg.Value)
		if err != nil {
			return err
		}

		if purchaseResult.Success {
			return a.createOrder(ctx, purchaseResult.Purchase)
		}

		return a.rollbackUpdateProductInventory(ctx, purchaseResult.Purchase)
	case common.RollbackProductInventoryHandler:
		purchaseResult, err := decodePbResponseToEventModel(msg.Value)
		if err != nil {
			return err
		}

		if !purchaseResult.Success {
			return a.publishPurchaseResult(ctx, &event.PurchaseResult{
				PurchaseID: purchaseResult.Purchase.ID,
				Status:     event.StatusRollbackFailed,
				Step:       event.StepUpdateProductInventory,
			})
		}

		return nil
	case common.CreateOrderHandler:
		purchaseResult, err := decodePbResponseToEventModel(msg.Value)
		if err != nil {
			return err
		}

		if purchaseResult.Success {
			// TODO: Implement create payment
			a.logger.Info("Creating payment...")
			return nil
		}

		return a.rollbackCreateOrder(ctx, purchaseResult.Purchase)
	case common.RollbackOrderHandler:
		purchaseResult, err := decodePbResponseToEventModel(msg.Value)
		if err != nil {
			return err
		}

		if !purchaseResult.Success {
			return a.publishPurchaseResult(ctx, &event.PurchaseResult{
				PurchaseID: purchaseResult.Purchase.ID,
				Status:     event.StatusRollbackFailed,
				Step:       event.StepCreateOrder,
			})
		}

		return nil
	default:
		return fmt.Errorf("handle reply: unknown handler: %s", msg.Headers[0].Key)
	}
}

func (a *app) rollbackUpdateProductInventory(ctx context.Context, purchase *aggregate.Purchase) error {
	err := a.publishPurchaseResult(ctx, &event.PurchaseResult{
		PurchaseID: purchase.ID,
		Status:     event.StatusFailed,
		Step:       event.StepUpdateProductInventory,
	})
	if err != nil {
		return err
	}

	pbRollback := &pb.RollbackPurchaseRequest{
		PurchaseId: purchase.ID,
		Timestamp:  timeconvert.Time2pbTimestamp(time.Now()),
	}

	err = a.publishPurchaseResult(ctx, &event.PurchaseResult{
		PurchaseID: purchase.ID,
		Status:     event.StatusRollback,
		Step:       event.StepUpdateProductInventory,
	})
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(pbRollback)
	return a.producer.PublishMessage(ctx, kafka.Message{
		Topic: common.RollbackProductInventoryTopic,
		Value: payload,
	})
}

func (a *app) createOrder(ctx context.Context, purchase *aggregate.Purchase) error {
	err := a.publishPurchaseResult(ctx, &event.PurchaseResult{
		PurchaseID: purchase.ID,
		Status:     event.StatusSucess,
		Step:       event.StepUpdateProductInventory,
	})
	if err != nil {
		return err
	}

	pbPurchase := encodeModel2PurchaseRequest(purchase)

	err = a.publishPurchaseResult(ctx, &event.PurchaseResult{
		PurchaseID: purchase.ID,
		Status:     event.StatusExecute,
		Step:       event.StepCreateOrder,
	})
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(pbPurchase)
	return a.producer.PublishMessage(ctx, kafka.Message{
		Topic: common.CreateOrderTopic,
		Value: payload,
	})
}

func (a *app) rollbackCreateOrder(ctx context.Context, purchase *aggregate.Purchase) error {
	err := a.publishPurchaseResult(ctx, &event.PurchaseResult{
		PurchaseID: purchase.ID,
		Status:     event.StatusFailed,
		Step:       event.StepCreateOrder,
	})
	if err != nil {
		return err
	}

	pbRollback := &pb.RollbackPurchaseRequest{
		PurchaseId: purchase.ID,
		Timestamp:  timeconvert.Time2pbTimestamp(time.Now()),
	}

	err = a.publishPurchaseResult(ctx, &event.PurchaseResult{
		PurchaseID: purchase.ID,
		Status:     event.StatusRollback,
		Step:       event.StepCreateOrder,
	})
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(pbRollback)
	err = a.producer.PublishMessage(ctx, kafka.Message{
		Topic: common.RollbackOrderTopic,
		Value: payload,
	})
	if err != nil {
		a.logger.Errorf("Orchestrator.RollbackFromOrder.RollbackCreateOrder, err: %v", err)
		return err
	}

	err = a.rollbackUpdateProductInventory(ctx, purchase)
	if err != nil {
		a.logger.Errorf("Orchestrator.RollbackFromOrder.RollbackUpdateProductInventory, err: %v", err)
	}
	return err
}

func (a *app) publishPurchaseResult(ctx context.Context, result *event.PurchaseResult) error {
	pbResult := encodePurchaseResult(result)

	payload, _ := json.Marshal(pbResult)
	return a.producer.PublishMessage(ctx, kafka.Message{
		Topic: common.PurchaseResultTopic,
		Value: payload,
	})
}
