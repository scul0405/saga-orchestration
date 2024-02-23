package eventhandler

import (
	"context"
	"encoding/json"
	"github.com/scul0405/saga-orchestration/internal/common"
	"github.com/scul0405/saga-orchestration/internal/purchase/domain/aggregate"
	kafkaClient "github.com/scul0405/saga-orchestration/pkg/kafka"
	"github.com/segmentio/kafka-go"
)

type PurchaseEventHandler interface {
	CreatePurchase(ctx context.Context, purchase *aggregate.Purchase) error
}

type purchaseEventHandler struct {
	producer kafkaClient.Producer
}

func NewPurchaseEventHandler(producer kafkaClient.Producer) PurchaseEventHandler {
	return &purchaseEventHandler{
		producer: producer,
	}
}

func (h *purchaseEventHandler) CreatePurchase(ctx context.Context, purchase *aggregate.Purchase) error {
	purchaseByte, _ := json.Marshal(purchase)
	msg := kafka.Message{
		Topic: common.PurchaseTopic,
		Key:   []byte("create_purchase"),
		Value: purchaseByte,
	}

	return h.producer.PublishMessage(ctx, msg)
}
