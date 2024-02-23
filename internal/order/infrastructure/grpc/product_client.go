package grpc

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/scul0405/saga-orchestration/internal/order/domain/valueobject"
	"github.com/scul0405/saga-orchestration/internal/pkg/grpcconn"
	pb "github.com/scul0405/saga-orchestration/proto"
)

type ProductService interface {
	GetProducts(ctx context.Context, productIDs *[]uint64) (*[]valueobject.DetailedPurchasedProduct, error)
}

type productServiceImpl struct {
	product endpoint.Endpoint
}

func NewProductService(conn *grpcconn.GRPCClientConn) ProductService {
	productSvc := grpcconn.NewGRPCClient("product.ProductService", "GetProducts", conn, &pb.GetProductsResponse{})

	return &productServiceImpl{
		product: productSvc,
	}
}

func (s *productServiceImpl) GetProducts(ctx context.Context, productIDs *[]uint64) (*[]valueobject.DetailedPurchasedProduct, error) {
	response, err := s.product(ctx, &pb.GetProductsRequest{ProductIds: *productIDs})
	if err != nil {
		return nil, err
	}

	pbProducts := response.(*pb.GetProductsResponse)

	products := make([]valueobject.DetailedPurchasedProduct, len(pbProducts.Products))

	for i, pbProduct := range pbProducts.Products {
		products[i] = valueobject.DetailedPurchasedProduct{
			ID:          pbProduct.Id,
			CategoryID:  pbProduct.Category,
			Name:        pbProduct.Name,
			BrandName:   pbProduct.BrandName,
			Description: pbProduct.Description,
			Price:       pbProduct.Price,
		}
	}

	return &products, nil
}
