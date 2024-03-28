package usecase

import (
	"context"
	"time"

	"github.com/Alladan04/vk_test/internal/models"
	"github.com/Alladan04/vk_test/internal/pkg/market"
	"github.com/satori/uuid"
)

type MarketUsecase struct {
	repo market.MarketRepo
}

func NewMarketUsecase(repo market.MarketRepo) *MarketUsecase {
	return &MarketUsecase{
		repo: repo,
	}
}

func (uc *MarketUsecase) AddItem(ctx context.Context, data models.ItemForm, userId uuid.UUID) (models.MarketItem, error) {
	err := data.Validate()
	if err != nil {
		return models.MarketItem{}, err
	}

	item := models.MarketItem{
		Id:          uuid.NewV4(),
		Title:       data.Title,
		Description: data.Description,
		ImagePath:   data.ImagePath,
		Price:       data.Price,
		Owner:       userId,
		CreateTime:  time.Now().UTC(),
	}
	err = uc.repo.AddItem(ctx, item)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (uc *MarketUsecase) GetAll(ctx context.Context, count int64, offset int64, sortOrder string, sortField string, minPrice float64, maxPrice float64) ([]models.MarketItem, error) {

	if minPrice == 0 && maxPrice == 0 {
		data, err := uc.repo.GetAll(ctx, count, offset, sortOrder, sortField)
		return data, err

	}
	data, err := uc.repo.GetFilteredByPrice(ctx, count, offset, sortOrder, sortField, minPrice, maxPrice)
	return data, err
}
