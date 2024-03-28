package market

import (
	"context"

	"github.com/Alladan04/vk_test/internal/models"
	"github.com/satori/uuid"
)

type MarketRepo interface {
	AddItem(context.Context, models.MarketItem) error
	GetAll(context.Context, int64, int64, string, string) ([]models.MarketItem, error)
	GetFilteredByPrice(context.Context, int64, int64, string, string, float64, float64) ([]models.MarketItem, error)
	GetMaxPrice(context.Context) (float64, error)
	GetMinPrice(context.Context) (float64, error)
}
type MarketUsecase interface {
	AddItem(context.Context, models.ItemForm, uuid.UUID) (models.MarketItem, error)
	GetAll(context.Context, int64, int64, string, string, string, string, uuid.UUID) ([]models.MarketItemResponse, error)
}
