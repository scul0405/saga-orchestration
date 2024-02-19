package command

import (
	"context"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
	"github.com/scul0405/saga-orchestration/services/product/internal/domain"
	"github.com/scul0405/saga-orchestration/services/product/internal/domain/entity"
	"github.com/scul0405/saga-orchestration/services/product/internal/domain/valueobject"
)

type CreateProduct struct {
	Name        string
	BrandName   string
	Description string
	Price       uint64
	CategoryID  uint64
	Inventory   uint64
}

type CreateProductHandler CommandHandler[CreateProduct]

type createProductHandler struct {
	sf          sonyflake.IDGenerator
	logger      logger.Logger
	productRepo domain.ProductRepository
}

func NewCreateProductHandler(sf sonyflake.IDGenerator, logger logger.Logger, productRepo domain.ProductRepository) CreateProductHandler {
	return &createProductHandler{
		sf:          sf,
		logger:      logger,
		productRepo: productRepo,
	}
}

func (h *createProductHandler) Handle(ctx context.Context, cmd CreateProduct) error {
	productID, err := h.sf.NextID()
	if err != nil {
		return err
	}

	err = h.productRepo.CreateProduct(ctx, &entity.Product{
		ID:         productID,
		CategoryID: cmd.CategoryID,
		Detail: &valueobject.ProductDetail{
			Name:        cmd.Name,
			Description: cmd.Description,
			BrandName:   cmd.BrandName,
			Price:       cmd.Price,
		},
		Inventory: cmd.Inventory,
	})

	if err != nil {
		return err
	}

	return nil
}
