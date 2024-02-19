package product

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/order/domain/valueobject"
	pb "github.com/scul0405/saga-orchestration/proto"
)

func encodeProductRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func decodeProductResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	return grpcReply, nil
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
