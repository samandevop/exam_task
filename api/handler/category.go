package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"crud/models"

	"github.com/gin-gonic/gin"
)

// CreateCategory godoc
// @ID create_category
// @Router /category [POST]
// @Summary Create Category
// @Description Create Category
// @Tags Category
// @Accept json
// @Produce json
// @Param category body models.CreateCategory true "CreateCategoryRequestBody"
// @Success 201 {object} models.Category "GetCategoryBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) CreateCategory(c *gin.Context) {
	var category models.CreateCategory

	err := c.ShouldBindJSON(&category)
	if err != nil {
		log.Printf("error whiling create: %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Category().Create(context.Background(), &category)
	if err != nil {
		log.Printf("error whiling Create: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling Create").Error())
		return
	}

	resp, err := h.storage.Category().GetByPKey(
		context.Background(),
		&models.CategoryPrimaryKey{Id: id},
	)

	if err != nil {
		log.Printf("error whiling GetByPKey: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GetByPKey").Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetByIdCategory godoc
// @ID get_by_id_category
// @Router /category/{id} [GET]
// @Summary Get By Id Category
// @Description Get By Id Category
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Category "GetCategoryBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) GetCategoryById(c *gin.Context) {

	id := c.Param("id")

	resp, err := h.storage.Category().GetByPKey(
		context.Background(),
		&models.CategoryPrimaryKey{Id: id},
	)

	if err != nil {
		log.Printf("error whiling GetByPKey: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GetByPKey").Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetListCategory godoc
// @ID get_list_category
// @Router /category [GET]
// @Summary Get List Category
// @Description Get List Category
// @Tags Category
// @Accept json
// @Produce json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} models.GetListCategoryResponse "GetCategoryBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) GetCategoryList(c *gin.Context) {
	var (
		limit  int
		offset int
		err    error
	)

	limitStr := c.Query("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("error whiling limit: %v\n", err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	offsetStr := c.Query("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			log.Printf("error whiling limit: %v\n", err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	resp, err := h.storage.Category().GetList(
		context.Background(),
		&models.GetListCategoryRequest{
			Limit:  int32(limit),
			Offset: int32(offset),
		},
	)

	if err != nil {
		log.Printf("error whiling get list: %v", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling get list").Error())
		return
	}

	c.JSON(http.StatusOK, resp)

}

// UpdateCategory godoc
// @ID update_category
// @Router /category/{id} [PUT]
// @Summary Update Category
// @Description Update Category
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param category body models.UpdateCategorySwagger true "CreateCategoryRequestBody"
// @Success 200 {object} models.Category "GetCategorysBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) UpdateCategory(c *gin.Context) {

	var (
		category models.UpdateCategory
	)

	id := c.Param("id")

	if id == "" {
		log.Printf("error whiling update: %v\n", errors.New("required category id").Error())
		c.JSON(http.StatusBadRequest, errors.New("required category id").Error())
		return
	}

	err := c.ShouldBindJSON(&category)
	if err != nil {
		log.Printf("error whiling update: %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	category.Id = id

	rowsAffected, err := h.storage.Category().Update(
		context.Background(),
		&category,
	)

	if err != nil {
		log.Printf("error whiling update: %v", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling update").Error())
		return
	}

	if rowsAffected == 0 {
		log.Printf("error whiling update rows affected: %v", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling update rows affected").Error())
		return
	}

	resp, err := h.storage.Category().GetByPKey(
		context.Background(),
		&models.CategoryPrimaryKey{Id: id},
	)

	if err != nil {
		log.Printf("error whiling GetByPKey: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GetByPKey").Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteByIdCategory godoc
// @ID delete_by_id_category
// @Router /category/{id} [DELETE]
// @Summary Delete By Id Category
// @Description Delete By Id Category
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Category "GetCategoryBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) DeleteCategory(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		log.Printf("error whiling update: %v\n", errors.New("required category id").Error())
		c.JSON(http.StatusBadRequest, errors.New("required category id").Error())
		return
	}

	err := h.storage.Category().Delete(
		context.Background(),
		&models.CategoryPrimaryKey{
			Id: id,
		},
	)

	if err != nil {
		log.Printf("error whiling delete: %v", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling delete").Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
