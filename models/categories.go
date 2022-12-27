package models

type CategoryPrimarKey struct {
	Id string `json:"id"`
}

type CreateCategory struct {
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
}

type Category struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ParentID  string `json:"parent_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type UpdateCategorySwagger struct {
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
}

type UpdateCategory struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
}

type GetListCategoryRequest struct {
	Limit  int32
	Offset int32
}

type GetListCategoryResponse struct {
	Count      int            `json:"count"`
	Categories []CategoryList `json:"categories"`
}

type CategoryList struct {
	Id        string      `json:"id"`
	Name      string      `json:"name"`
	ParentID  string      `json:"parent_id"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	Childs    []ListChild `json:"childs"`
}

type ListChild struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ParentID  string `json:"parent_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
