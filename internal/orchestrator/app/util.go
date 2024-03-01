package app

import (
	"encoding/json"
	"github.com/scul0405/saga-orchestration/internal/orchestrator/domain/aggregate"
	"github.com/scul0405/saga-orchestration/internal/orchestrator/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/orchestrator/domain/event"
	"github.com/scul0405/saga-orchestration/internal/orchestrator/domain/valueobject"
	"github.com/scul0405/saga-orchestration/pkg/timeconvert"
	pb "github.com/scul0405/saga-orchestration/proto"
	"time"
)

func encodePurchaseResult(result *event.PurchaseResult) *pb.PurchaseResult {
	return &pb.PurchaseResult{
		PurchaseId: result.PurchaseID,
		Step:       getPbPurchaseStep(result.Step),
		Status:     getPbPurchaseStatus(result.Status),
		Timestamp:  timeconvert.Time2pbTimestamp(time.Now()),
	}
}

func encodeModel2PurchaseRequest(purchase *aggregate.Purchase) *pb.CreatePurchaseRequest {
	orderItems := make([]*pb.PurchaseOrderItem, len(*(purchase.Order.OrderItems)))
	for i, item := range *purchase.Order.OrderItems {
		orderItems[i] = &pb.PurchaseOrderItem{
			ProductId: item.ID,
			Quantity:  item.Quantity,
		}
	}

	return &pb.CreatePurchaseRequest{
		PurchaseId: purchase.ID,
		Purchase: &pb.Purchase{
			Order: &pb.Order{
				CustomerId: purchase.Order.CustomerID,
				OrderItems: orderItems,
			},
			Payment: &pb.Payment{
				CurrencyCode: purchase.Payment.CurrencyCode,
				Amount:       purchase.Payment.Amount,
			},
		},
	}
}

func decodePbResponseToEventModel(data []byte) (*event.CreatePurchaseResponse, error) {
	var pbResult pb.CreatePurchaseResponse
	err := json.Unmarshal(data, &pbResult)
	if err != nil {
		return nil, err
	}

	orderItems := make([]entity.OrderItem, len(pbResult.Purchase.Order.OrderItems))
	for i, item := range pbResult.Purchase.Order.OrderItems {
		orderItems[i] = entity.OrderItem{
			ID:       item.ProductId,
			Quantity: item.Quantity,
		}
	}

	return &event.CreatePurchaseResponse{
		Purchase: &aggregate.Purchase{
			ID: pbResult.PurchaseId,
			Order: &entity.Order{
				CustomerID: pbResult.Purchase.Order.CustomerId,
				OrderItems: &orderItems,
			},
			Payment: &valueobject.Payment{
				CurrencyCode: pbResult.Purchase.Payment.CurrencyCode,
				Amount:       pbResult.Purchase.Payment.Amount,
			},
		},
		Success: pbResult.Success,
		Error:   pbResult.ErrorMessage,
	}, nil
}

func getPbPurchaseStep(step string) pb.PurchaseStep {
	switch step {
	case event.StepUpdateProductInventory:
		return pb.PurchaseStep_UPDATE_PRODUCT_INVENTORY
	case event.StepCreateOrder:
		return pb.PurchaseStep_CREATE_ORDER
	case event.StepCreatePayment:
		return pb.PurchaseStep_CREATE_PAYMENT
	}
	return -1
}

func getPbPurchaseStatus(status string) pb.PurchaseStatus {
	switch status {
	case event.StatusExecute:
		return pb.PurchaseStatus_EXECUTE
	case event.StatusSucess:
		return pb.PurchaseStatus_SUCCESS
	case event.StatusFailed:
		return pb.PurchaseStatus_FAILED
	case event.StatusRollbacked:
		return pb.PurchaseStatus_ROLLBACKED
	case event.StatusRollbackFailed:
		return pb.PurchaseStatus_ROLLBACK_FAILED
	}
	return -1
}
