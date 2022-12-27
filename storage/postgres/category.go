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
		id     = uuid.New().String()
		query  string
		nullId sql.NullString
	)

	query = `
		INSERT INTO categories(
			id,
			name,
			parent_id,
			updated_at
		) VALUES ( $1, $2, $3, now() )
	`

	if category.ParentID == "" {
		_, err := f.db.Exec(ctx, query,
			id,
			category.Name,
			nullId,
		)

		if err != nil {
			return "", err
		}
	} else {
		_, err := f.db.Exec(ctx, query,
			id,
			category.Name,
			category.ParentID,
		)

		if err != nil {
			return "", err
		}
	}

	return id, nil
}

func (f *CategoryRepo) GetByPKey(ctx context.Context, pkey *models.CategoryPrimarKey) (*models.CategoryList, error) {

	var (
		resp = models.CategoryList{}
	)

	query := `
	SELECT
		c1.id,
		c1.name,
		ARRAY_AGG(c2.id),
		ARRAY_AGG(c2.name),
		ARRAY_AGG(c2.parent_id)
	FROM
		categories as c1
	JOIN categories as c2 ON c1.id = c2.parent_id
	WHERE c1.is_deleted = false AND c2.is_deleted = false AND c1.id = $1
	GROUP BY c1.id, c1.name
	`

	rows, err := f.db.Query(ctx, query, pkey.Id)

	for rows.Next() {
		var (
			respChild = []models.ListChild{}
			id        sql.NullString
			name      sql.NullString
			childId   []string
			childName []string
			childPI   []string
		)

		err := rows.Scan(
			&id,
			&name,
			&childId,
			&childName,
			&childPI,
		)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(childId); i++ {
			respChild = append(respChild, models.ListChild{
				Id:       childId[i],
				Name:     childName[i],
				ParentID: childPI[i],
			})

		}

		resp.Id = id.String
		resp.Name = name.String
		resp.Childs = respChild

	}

	return &resp, err
}

func (f *CategoryRepo) GetList(ctx context.Context, req *models.GetListCategoryRequest) (*models.GetListCategoryResponse, error) {

	var (
		resp   = models.GetListCategoryResponse{}
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
		c1.id,
		c1.name,
		ARRAY_AGG(c2.id),
		ARRAY_AGG(c2.name),
		ARRAY_AGG(c2.parent_id)
	FROM
		categories as c1
	JOIN categories as c2 ON c1.id = c2.parent_id
	WHERE c1.is_deleted = false AND c2.is_deleted = false
	GROUP BY c1.id, c1.name
	`

	query += offset + limit

	rows, err := f.db.Query(ctx, query)

	for rows.Next() {
		var (
			respChild = []models.ListChild{}
			id        sql.NullString
			name      sql.NullString
			childId   []string
			childName []string
			childPI   []string
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&childId,
			&childName,
			&childPI,
		)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(childId); i++ {
			respChild = append(respChild, models.ListChild{
				Id:       childId[i],
				Name:     childName[i],
				ParentID: childPI[i],
			})

		}

		resp.Categories = append(resp.Categories, models.CategoryList{
			Id:     id.String,
			Name:   name.String,
			Childs: respChild,
		})

	}

	return &resp, err
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
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":        req.Id,
		"name":      req.Name,
		"parent_id": req.ParentID,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	rowsAffected, err := f.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (f *CategoryRepo) Delete(ctx context.Context, req *models.CategoryPrimarKey) error {

	_, err := f.db.Exec(ctx, "UPDATE categories SET deleted_at = now(), is_deleted = true WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return err
}
