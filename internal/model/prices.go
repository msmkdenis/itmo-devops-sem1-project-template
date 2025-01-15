package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type LoadResult struct {
	TotalQuantity   int             `db:"total_items" json:"total_items"` //nolint:tagliatelle //неправильно обрабатывает
	TotalCategories int             `db:"total_categories" json:"total_categories"`
	TotalPrice      decimal.Decimal `db:"total_price" json:"total_price"`
}

type Product struct {
	ID         int             `db:"id"`
	Name       string          `db:"name"`
	Category   string          `db:"category"`
	Price      decimal.Decimal `db:"price"`
	CreateDate time.Time       `db:"create_date"`
}
