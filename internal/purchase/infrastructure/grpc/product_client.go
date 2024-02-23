package grpc

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/scul0405/saga-orchestration/internal/pkg/grpcconn"
	"github.com/scul0405/saga-orchestration/internal/purchase/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/purchase/domain/valueobject"
	pb "github.com/scul0405/saga-orchestration/proto"
)

type ProductService interface {
	CheckProducts(ctx context.Context, orderItems *[]entity.OrderItem) (*[]valueobject.ProductStatus, error)
}

type productServiceImpl struct {
	product endpoint.Endpoint
}

func NewProductService(conn *grpcconn.GRPCClientConn) ProductService {
	productSvc := grpcconn.NewGRPCClient("product.ProductService", "CheckProducts", conn, &pb.CheckProductsResponse{})

	return &productServiceImpl{
		product: productSvc,
	}
}

func (s *productServiceImpl) CheckProducts(ctx context.Context, orderItems *[]entity.OrderItem) (*[]valueobject.ProductStatus, error) {
	pbOrderItems := make([]*pb.OrderItem, len(*orderItems))
	for i, item := range *orderItems {
		pbOrderItems[i] = &pb.OrderItem{
			ProductId: item.ID,
			Quantity:  item.Quantity,
		}
	}

	response, err := s.product(ctx, &pb.CheckProductsRequest{Items: pbOrderItems})
	if err != nil {
		return nil, err
	}

	pbProductStatuses := response.(*pb.CheckProductsResponse)

	productStatuses := make([]valueobject.ProductStatus, len(pbProductStatuses.Statuses))
	for i, pbProductStatus := range pbProductStatuses.Statuses {
		productStatuses[i] = valueobject.ProductStatus{
			ProductID: pbProductStatus.ProductId,
			Status:    getProductStatus(pbProductStatus.Status),
			Price:     pbProductStatus.Price,
		}
	}

	return &productStatuses, nil
}

func getProductStatus(status pb.Status) valueobject.Status {
	switch status {
	case pb.Status_OK:
		return valueobject.ProductOk
	case pb.Status_NOT_FOUND:
		return valueobject.ProductNotFound
	default:
		return -1
	}
}
