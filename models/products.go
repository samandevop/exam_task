package models

type ProductPrimarKey struct {
	Id string `json:"id"`
}

type CreateProduct struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	CategoryID string  `json:"category_id"`
}

type Product struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	CategoryID string  `json:"category_id"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	DeletedAt  string  `json:"deleted_at"`
}

type UpdateProductSwagger struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	CategoryID string  `json:"category_id"`
}

type UpdateProduct struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	CategoryID string  `json:"category_id"`
}

type GetListProductRequest struct {
	Limit  int32
	Offset int32
}

type GetListProductResponse struct {
	Count    int       `json:"count"`
	Products []Product `json:"products"`
}
