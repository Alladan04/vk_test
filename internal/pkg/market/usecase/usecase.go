package usecase

import (
	"context"
	"errors"
	"strconv"
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

func (uc *MarketUsecase) GetAll(ctx context.Context, count int64, offset int64, sortOrder string, sortField string, minParam string, maxParam string, userId uuid.UUID) ([]models.MarketItemResponse, error) {
	var minPrice float64
	var maxPrice float64
	//find out min and max real values
	if minParam == "" {
		minPrice, _ = uc.repo.GetMinPrice(ctx)

	} else {
		var err error
		minPrice, err = strconv.ParseFloat(minParam, 64)
		if err != nil {
			return nil, errors.New("wrong minPrice param")
		}
	}
	if maxParam == "" {
		maxPrice, _ = uc.repo.GetMaxPrice(ctx)
	} else {
		var err error
		maxPrice, err = strconv.ParseFloat(maxParam, 64)
		if err != nil {
			return nil, errors.New("wrong maxPrice param")
		}
	}

	data, err := uc.repo.GetFilteredByPrice(ctx, count, offset, sortOrder, sortField, minPrice, maxPrice)

	if err != nil {
		return nil, err
	}
	result := make([]models.MarketItemResponse, 0)

	for _, element := range data {
		result = append(result,
			models.MarketItemResponse{
				Item:           element,
				IsCurrentUsers: element.Owner == userId},
		)

	}
	return result, err
}
