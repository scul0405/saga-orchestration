package command

import (
	"context"

	"github.com/scul0405/saga-orchestration/internal/product/domain"
	"github.com/scul0405/saga-orchestration/internal/product/domain/entity"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
)

type CreateCategory struct {
	Name        string
	Description string
}

type CreateCategoryHandler CommandHandler[CreateCategory]

type createCategoryHandler struct {
	sf           sonyflake.IDGenerator
	logger       logger.Logger
	categoryRepo domain.CategoryRepository
}

func NewCreateCategoryHandler(sf sonyflake.IDGenerator, logger logger.Logger, categoryRepo domain.CategoryRepository) CreateCategoryHandler {
	return &createCategoryHandler{
		sf:           sf,
		logger:       logger,
		categoryRepo: categoryRepo,
	}
}

func (h *createCategoryHandler) Handle(ctx context.Context, cmd CreateCategory) error {
	productID, err := h.sf.NextID()
	if err != nil {
		return err
	}

	err = h.categoryRepo.CreateCategory(ctx, &entity.Category{
		ID:          productID,
		Name:        cmd.Name,
		Description: cmd.Description,
	})

	if err != nil {
		return err
	}

	return nil
}
