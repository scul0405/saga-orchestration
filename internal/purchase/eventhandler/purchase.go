package eventhandler

import (
	"context"
	"encoding/json"
	"github.com/scul0405/saga-orchestration/internal/common"
	kafkaClient "github.com/scul0405/saga-orchestration/pkg/kafka"
	pb "github.com/scul0405/saga-orchestration/proto"
	"github.com/segmentio/kafka-go"
)

type PurchaseEventHandler interface {
	ProduceCreatePurchase(ctx context.Context, purchase *pb.CreatePurchaseRequest) error
}

type purchaseEventHandler struct {
	producer kafkaClient.Producer
}

func NewPurchaseEventHandler(producer kafkaClient.Producer) PurchaseEventHandler {
	return &purchaseEventHandler{
		producer: producer,
	}
}

func (h *purchaseEventHandler) ProduceCreatePurchase(ctx context.Context, purchase *pb.CreatePurchaseRequest) error {
	purchaseByte, _ := json.Marshal(purchase)
	msg := kafka.Message{
		Topic: common.PurchaseTopic,
		Key:   []byte("create_purchase"),
		Value: purchaseByte,
	}

	return h.producer.PublishMessage(ctx, msg)
}
