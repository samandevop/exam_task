package storage

import (
	"context"

	"crud/models"
)

type StorageI interface {
	CloseDB()
	Category() CategoryRepoI
	Product() ProductRepoI
	Order() OrderRepoI
}

type CategoryRepoI interface {
	Create(ctx context.Context, req *models.CreateCategory) (string, error)
	GetByPKey(ctx context.Context, req *models.CategoryPrimaryKey) (*models.CategoryList, error)
	GetList(ctx context.Context, req *models.GetListCategoryRequest) (*models.GetListCategoryResponse, error)
	Update(ctx context.Context, req *models.UpdateCategory) (int64, error)
	Delete(ctx context.Context, req *models.CategoryPrimaryKey) error
}

type ProductRepoI interface {
	Create(ctx context.Context, req *models.CreateProduct) (string, error)
	GetByPKey(ctx context.Context, req *models.ProductPrimarKey) (*models.Product, error)
	GetList(ctx context.Context, req *models.GetListProductRequest) (*models.GetListProductResponse, error)
	Update(ctx context.Context, req *models.UpdateProduct) (int64, error)
	Delete(ctx context.Context, req *models.ProductPrimarKey) error
}

type OrderRepoI interface {
	Create(ctx context.Context, req *models.CreateOrder) (string, error)
	GetByPKey(ctx context.Context, req *models.OrderPrimarKey) (*models.OrderList, error)
	GetList(ctx context.Context, req *models.GetListOrderRequest) (*models.GetListOrderResponse, error)
	Update(ctx context.Context, req *models.UpdateOrder) (int64, error)
	Delete(ctx context.Context, req *models.OrderPrimarKey) error
}
