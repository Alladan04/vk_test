package repo

import (
	"context"
	"fmt"

	"github.com/Alladan04/vk_test/internal/models"
	"github.com/jackc/pgtype/pgxtype"
)

const (
	addItem            = "INSERT INTO market (id, title, description, image_path, price, create_time, owner) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	getAll             = "SELECT id, title, description, image_path, price, create_time, owner FROM market ORDER BY %s %s LIMIT $1 OFFSET $2; "
	getFilteredByPrice = "SELECT id, title, description, image_path, price, create_time, owner FROM market WHERE price>$1 AND price<$2 ORDER BY %s %s LIMIT $3 OFFSET $4;"
	getMinPrice        = "SELECT MIN(price) FROM market;"
	getMaxPrice        = "SELECT MAX(price) FROM market;"
)

type MarketRepo struct {
	db pgxtype.Querier
}

func NewMarketRepo(db pgxtype.Querier) *MarketRepo {
	return &MarketRepo{
		db: db,
	}
}

func (repo *MarketRepo) AddItem(ctx context.Context, item models.MarketItem) error {
	_, err := repo.db.Exec(ctx, addItem, item.Id, item.Title, item.Description, item.ImagePath, item.Price, item.CreateTime, item.Owner)
	if err != nil {
		return err
	}
	return nil

}
func (repo *MarketRepo) GetAll(ctx context.Context, count int64, offset int64, sortOrder string, sortField string) ([]models.MarketItem, error) {
	result := make([]models.MarketItem, 0, count)
	query := fmt.Sprintf(getAll, sortField, sortOrder)
	fmt.Println(query)
	rows, err := repo.db.Query(ctx, query, count, offset)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var item models.MarketItem
		if err := rows.Scan(&item.Id, &item.Title, &item.Description, &item.ImagePath, &item.Price, &item.CreateTime, &item.Owner); err != nil {
			return result, fmt.Errorf("error occured while scanning items:%w", err)
		}
		result = append(result, item)
	}

	return result, nil
}

func (repo *MarketRepo) GetFilteredByPrice(ctx context.Context, count int64, offset int64, sortOrder string, sortField string, minPrice float64, maxPrice float64) ([]models.MarketItem, error) {
	result := make([]models.MarketItem, 0, count)
	query := fmt.Sprintf(getFilteredByPrice, sortField, sortOrder)
	rows, err := repo.db.Query(ctx, query, minPrice, maxPrice, count, offset)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var item models.MarketItem
		if err := rows.Scan(&item.Id, &item.Title, &item.Description, &item.ImagePath, &item.Price, &item.CreateTime, &item.Owner); err != nil {
			return result, fmt.Errorf("error occured while scanning items:%w", err)
		}
		result = append(result, item)
	}

	return result, nil
}
