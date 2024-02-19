package product

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/scul0405/saga-orchestration/cmd/order/config"
	"github.com/scul0405/saga-orchestration/internal/order/domain/valueobject"
	pb "github.com/scul0405/saga-orchestration/proto"
	"github.com/sony/gobreaker"
	"time"
)

const (
	ProductSvcName = "product.ProductService"
	timeout        = 5 * time.Second
)

type ProductService interface {
	GetProducts(ctx context.Context, productIDs *[]uint64) (*[]valueobject.DetailedPurchasedProduct, error)
}

type productServiceImpl struct {
	product endpoint.Endpoint
}

func NewProductService(cfg *config.Config, conn *ProductConn) ProductService {

	var opts []grpc.ClientOption

	var productEndpoint endpoint.Endpoint

	{
		productEndpoint = grpc.NewClient(
			conn.Conn,
			ProductSvcName,
			"GetProducts",
			encodeProductRequest,
			decodeProductResponse,
			&pb.GetProductsResponse{},
			append(opts, grpc.ClientBefore(grpc.SetRequestHeader("Service-Name", ProductSvcName)))...,
		).Endpoint()

		productEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    ProductSvcName,
			Timeout: timeout,
		}))(productEndpoint)
	}

	return &productServiceImpl{
		product: productEndpoint,
	}
}
