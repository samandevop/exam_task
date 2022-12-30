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

type CategoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{
		db: db,
	}
}

func (f *CategoryRepo) Create(ctx context.Context, category *models.CreateCategory) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO categories (
			id,
			name,
			parent_id,
			updated_at
		) VALUES ( $1, $2, $3, now() )
	`

	_, err := f.db.Exec(ctx, query,
		id,
		category.Name,
		helper.NewNullString(category.ParentID),
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (f *CategoryRepo) GetByPKey(ctx context.Context, pkey *models.CategoryPrimaryKey) (*models.CategoryList, error) {

	var (
		id        sql.NullString
		name      sql.NullString
		parentID  sql.NullString
		createdAt sql.NullString
		updatedAt sql.NullString
	)

	query := `
		SELECT
			id,
			name,
			parent_id,
			created_at,
			updated_at
		FROM categories
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := f.db.QueryRow(ctx, query, pkey.Id).Scan(
		&id,
		&name,
		&parentID,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	resp := &models.CategoryList{
		Id:        id.String,
		Name:      name.String,
		ParentID:  parentID.String,
		CreatedAt: createdAt.String,
		UpdatedAt: updatedAt.String,
	}

	queryChild := `
		SELECT
			id,
			name,
			parent_id,
			created_at,
			updated_at
		FROM categories
		WHERE parent_id = $1 AND deleted_at IS NULL
	`

	rows, err := f.db.Query(ctx, queryChild, resp.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, nil
		}

		return nil, err
	}

	for rows.Next() {

		err = rows.Scan(
			&id,
			&name,
			&parentID,
			&createdAt,
			&updatedAt,
		)

		resp.Childs = append(resp.Childs, &models.Category{
			Id:        id.String,
			Name:      name.String,
			ParentID:  parentID.String,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})
	}

	return resp, err
}

func (f *CategoryRepo) GetList(ctx context.Context, req *models.GetListCategoryRequest) (*models.GetListCategoryResponse, error) {

	var (
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		resp   = &models.GetListCategoryResponse{}
	)

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query := `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			parent_id,
			created_at,
			updated_at
		FROM categories
		WHERE parent_id IS NULL AND deleted_at IS NULL
	`

	query += offset + limit

	rows, err := f.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	// get parent categories
	for rows.Next() {

		var (
			id        sql.NullString
			name      sql.NullString
			parentID  sql.NullString
			createdAt sql.NullString
			updatedAt sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&name,
			&parentID,
			&createdAt,
			&updatedAt,
		)

		resp.Categories = append(resp.Categories, &models.CategoryList{
			Id:        id.String,
			Name:      name.String,
			ParentID:  parentID.String,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})
	}

	// get category childs
	for _, category := range resp.Categories {

		queryChild := `
			SELECT
				id,
				name,
				parent_id,
				created_at,
				updated_at
			FROM categories
			WHERE parent_id = $1 AND deleted_at IS NULL
		`

		rows, err := f.db.Query(ctx, queryChild, category.Id)
		if err != nil {
			if err == sql.ErrNoRows {
				return resp, nil
			}

			return nil, err
		}

		for rows.Next() {
			var (
				id        sql.NullString
				name      sql.NullString
				parentID  sql.NullString
				createdAt sql.NullString
				updatedAt sql.NullString
			)

			err = rows.Scan(
				&id,
				&name,
				&parentID,
				&createdAt,
				&updatedAt,
			)

			category.Childs = append(category.Childs, &models.Category{
				Id:        id.String,
				Name:      name.String,
				ParentID:  parentID.String,
				CreatedAt: createdAt.String,
				UpdatedAt: updatedAt.String,
			})
		}
	}

	return resp, err
}

func (f *CategoryRepo) Update(ctx context.Context, req *models.UpdateCategory) (int64, error) {

	var (
		query  = ""
		params map[string]interface{}
	)

	query = `
		UPDATE
			categories
		SET
			name = :name,
			parent_id = :parent_id,
			updated_at = now()
		WHERE id = :id AND deleted_at IS NULL
	`

	params = map[string]interface{}{
		"id":        req.Id,
		"name":      req.Name,
		"parent_id": helper.NewNullString(req.ParentID),
	}

	query, args := helper.ReplaceQueryParams(query, params)

	rowsAffected, err := f.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (f *CategoryRepo) Delete(ctx context.Context, req *models.CategoryPrimaryKey) error {

	_, err := f.db.Exec(ctx, "UPDATE categories SET deleted_at = now() WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return err
}
