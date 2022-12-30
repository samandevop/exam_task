package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"crud/models"
	"crud/pkg/helper"
)

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

func (f *ProductRepo) Create(ctx context.Context, product *models.CreateProduct) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO products(
			id,
			name,
			price,
			category_id,
			updated_at
		) VALUES ( $1, $2, $3, $4, now() )
	`

	_, err := f.db.Exec(ctx, query,
		id,
		product.Name,
		product.Price,
		product.CategoryID,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (f *ProductRepo) GetByPKey(ctx context.Context, pkey *models.ProductPrimarKey) (*models.Product, error) {

	var (
		id          sql.NullString
		name        sql.NullString
		price       sql.NullFloat64
		category_id sql.NullString
		createdAt   sql.NullString
		updatedAt   sql.NullString
	)

	query := `
		SELECT
			id,
			name,
			price,
			category_id,
			created_at,
			updated_at
		FROM
			products
		WHERE products.deleted_at IS NULL AND id = $1
	`

	err := f.db.QueryRow(ctx, query, pkey.Id).
		Scan(
			&id,
			&name,
			&price,
			&category_id,
			&createdAt,
			&updatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &models.Product{
		Id:         id.String,
		Name:       name.String,
		Price:      price.Float64,
		CategoryID: category_id.String,
		CreatedAt:  createdAt.String,
		UpdatedAt:  updatedAt.String,
	}, nil
}

func (f *ProductRepo) GetList(ctx context.Context, req *models.GetListProductRequest) (*models.GetListProductResponse, error) {

	var (
		resp   = models.GetListProductResponse{}
		offset = ""
		limit  = ""
	)

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	query := `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			price,
			category_id,
			created_at,
			updated_at
		FROM
			products
		WHERE products.deleted_at IS NULL
	`

	query += offset + limit

	rows, err := f.db.Query(ctx, query)

	for rows.Next() {

		var (
			id          sql.NullString
			name        sql.NullString
			price       sql.NullFloat64
			category_id sql.NullString
			createdAt   sql.NullString
			updatedAt   sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&price,
			&category_id,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Products = append(resp.Products, models.Product{
			Id:         id.String,
			Name:       name.String,
			Price:      price.Float64,
			CategoryID: category_id.String,
			CreatedAt:  createdAt.String,
			UpdatedAt:  updatedAt.String,
		})

	}

	return &resp, err
}

func (f *ProductRepo) Update(ctx context.Context, req *models.UpdateProduct) (int64, error) {

	var (
		query  = ""
		params map[string]interface{}
	)

	query = `
		UPDATE
			products
		SET
			name = :name,
			price = :price,
			category_id = :category_id,
			updated_at = now()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":          req.Id,
		"name":        req.Name,
		"price":       req.Price,
		"category_id": req.CategoryID,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	rowsAffected, err := f.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (f *ProductRepo) Delete(ctx context.Context, req *models.ProductPrimarKey) error {

	_, err := f.db.Exec(ctx, "UPDATE products SET deleted_at = now() WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return err
}
