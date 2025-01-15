package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"project_sem/internal/model"
	"project_sem/internal/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

//nolint:gochecknoglobals // Настройка, переиспользуемая в рамках пакета
var psql = getQueryBuilder()

type MarketingRepository struct {
	postgresPool *storage.PostgresPool
}

func NewMarketingRepository(postgresPool *storage.PostgresPool) *MarketingRepository {
	return &MarketingRepository{postgresPool: postgresPool}
}

func (r *MarketingRepository) UploadProducts(ctx context.Context, products []model.Product) (*model.LoadResult, error) {
	tx, err := r.postgresPool.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to acquire connection for transaction %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("failed to rollback transaction", "error", err)
		}
	}()

	insertProductQuery := `
							insert into "prices" (id, name, category, price, create_date) 
							values ($1, $2, $3, $4, $5) 
							on conflict (id) do nothing`

	productStatement, err := tx.Prepare(ctx, "insertproduct", insertProductQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare query %w", err)
	}

	batch := &pgx.Batch{}

	// технически for range имеет небольшие накладные расходы на копирование элемента, поэтому так сделал
	for i := 0; i < len(products); i++ {
		batch.Queue(productStatement.Name, products[i].ID, products[i].Name, products[i].Category, products[i].Price, products[i].CreateDate)
	}

	result := tx.SendBatch(ctx, batch)

	if err := result.Close(); err != nil {
		return nil, fmt.Errorf("error executing batch: %w", err)
	}

	query := `
		select
			count(id) as total_items,
			count(distinct category) as total_categories,
			coalesce(sum(price), 0) as total_price
		from prices`

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get products stats: %w", err)
	}

	defer rows.Close()

	stats, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.LoadResult])
	if err != nil {
		return nil, fmt.Errorf("collect product rows: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("unable to commit: %w", err)
	}

	return &stats, nil
}

func (r *MarketingRepository) LoadProducts(ctx context.Context) ([]model.Product, error) {
	queryBuilder := psql.Select(
		"id",
		"name",
		"category",
		"price",
		"create_date",
	).From("prices")

	query, _, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.postgresPool.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get products list: %w", err)
	}

	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Product])
	if err != nil {
		return nil, fmt.Errorf("collect product rows: %w", err)
	}

	return products, nil
}

func getQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
