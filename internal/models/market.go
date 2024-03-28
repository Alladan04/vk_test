package models

import (
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/satori/uuid"
)

const (
	titleMinLength       = 2
	titleMaxLength       = 100
	desscriptionMaxLenth = 3000
)

type MarketItem struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImagePath   string    `json:"image_path"`
	Price       float64   `json:"price"`
	CreateTime  time.Time `json:"create_time"`
	Owner       uuid.UUID `json:"owner"`
}

type MarketItemResponse struct {
	Item           MarketItem `json:"item"`
	IsCurrentUsers bool       `json:"is_current_users"`
}
type ItemForm struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImagePath   string  `json:"image_path,omitempty"`
	Price       float64 `json:"price"`
}

func (form *ItemForm) Validate() error {
	runedTitle := []rune(form.Title)
	runedDescription := []rune(form.Description)
	if len(runedTitle) < titleMinLength {
		return fmt.Errorf("title mast contain at least %d symbols", titleMinLength)

	}
	if len(runedTitle) > titleMaxLength {
		return fmt.Errorf("title mastn`t contain more than %d symbols", titleMaxLength)
	}
	if len(runedDescription) > desscriptionMaxLenth {
		return fmt.Errorf("description mastn`t contain more than %d symbols", desscriptionMaxLenth)
	}

	for _, sym := range runedTitle {
		if !unicode.IsDigit(sym) && !isEnglishLetter(sym) && sym != '-' {
			return errors.New("item Title can only include symbols: A-Z, a-z, 0-9")
		}
	}
	return nil
}
