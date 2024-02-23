package grpc

import (
	"context"
	"fmt"
	"github.com/scul0405/saga-orchestration/internal/product/app/query"
	"github.com/scul0405/saga-orchestration/internal/product/domain/valueobject"
	pb "github.com/scul0405/saga-orchestration/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *Server) CheckProducts(ctx context.Context, req *pb.CheckProductsRequest) (*pb.CheckProductsResponse, error) {
	items := req.GetItems()
	productIds := make([]uint64, len(items))

	for i, item := range items {
		productIds[i] = item.GetProductId()
	}

	productStatuses, err := srv.app.Queries.CheckProducts.Handle(ctx, query.CheckProducts{
		ProductIDs: &productIds,
	})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("internal error: %v", err),
		)
	}

	pbProductStatuses := make([]*pb.ProductStatus, len(*productStatuses))

	for i, productStatus := range *productStatuses {
		pbProductStatuses[i] = &pb.ProductStatus{
			ProductId: productStatus.ID,
			Status:    getPBStatus(productStatus),
			Price:     productStatus.Price,
		}
	}

	return &pb.CheckProductsResponse{
		Statuses: pbProductStatuses,
	}, nil
}

func (srv *Server) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	productIDs := req.GetProductIds()
	products, err := srv.app.Queries.GetProducts.Handle(ctx, query.GetProducts{
		ProductIDs: &productIDs,
	})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("internal error: %v", err),
		)
	}

	pbProducts := make([]*pb.Product, len(*products))

	for i, product := range *products {
		pbProducts[i] = &pb.Product{
			Id:          product.ID,
			Category:    product.CategoryID,
			Name:        product.Detail.Name,
			BrandName:   product.Detail.BrandName,
			Description: product.Detail.Description,
			Price:       product.Detail.Price,
			Inventory:   product.Inventory,
		}
	}

	return &pb.GetProductsResponse{
		Products: pbProducts,
	}, nil
}

func getPBStatus(productStatus valueobject.ProductStatus) pb.Status {
	if productStatus.Status {
		return pb.Status_OK
	}

	if productStatus.Price == 0 {
		return pb.Status_NOT_FOUND
	}

	return pb.Status_NOT_ENOUGH
}
