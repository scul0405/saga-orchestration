package command

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/purchase/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/purchase/eventhandler"
	"github.com/scul0405/saga-orchestration/internal/purchase/infrastructure/grpc"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
	pb "github.com/scul0405/saga-orchestration/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type CreatePurchase struct {
	Order   *Order
	Payment *Payment
}

type Order struct {
	CustomerID uint64
	OrderItems *[]OrderItem
}

type OrderItem struct {
	ID       uint64
	Price    uint64
	Quantity uint64
}

type Payment struct {
	CurrencyCode string
	Amount       uint64
}

type CreatePurchaseHandler CommandHandler[CreatePurchase]

type createPurchaseHandler struct {
	sf         sonyflake.IDGenerator
	logger     logger.Logger
	productSvc grpc.ProductService
	evPub      eventhandler.PurchaseEventHandler
}

func NewCreatePurchaseHandler(
	sf sonyflake.IDGenerator, logger logger.Logger,
	productSvc grpc.ProductService,
	evPub eventhandler.PurchaseEventHandler) CreatePurchaseHandler {
	return &createPurchaseHandler{
		sf:         sf,
		logger:     logger,
		productSvc: productSvc,
		evPub:      evPub,
	}
}

func (h *createPurchaseHandler) Handle(ctx context.Context, cmd CreatePurchase) error {
	orderItems := make([]entity.OrderItem, len(*cmd.Order.OrderItems))
	for i, item := range *cmd.Order.OrderItems {
		orderItems[i] = entity.OrderItem{
			ID:       item.ID,
			Quantity: item.Quantity,
		}
	}

	purchaseID, err := h.sf.NextID()
	if err != nil {
		return err
	}

	pbPurchaseOrderItem := make([]*pb.PurchaseOrderItem, len(*cmd.Order.OrderItems))
	for i, item := range *cmd.Order.OrderItems {
		pbPurchaseOrderItem[i] = &pb.PurchaseOrderItem{
			ProductId: item.ID,
			Quantity:  item.Quantity,
		}
	}

	purchase := &pb.CreatePurchaseRequest{
		PurchaseId: purchaseID,
		Purchase: &pb.Purchase{
			Order: &pb.Order{
				CustomerId: cmd.Order.CustomerID,
				OrderItems: pbPurchaseOrderItem,
			},
			Payment: &pb.Payment{
				CurrencyCode: cmd.Payment.CurrencyCode,
				Amount:       cmd.Payment.Amount,
			},
		},
		Timestamp: timestamppb.New(time.Now()),
	}

	return h.evPub.ProduceCreatePurchase(ctx, purchase)
}
