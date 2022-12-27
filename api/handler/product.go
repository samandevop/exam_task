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

// CreateProduct godoc
// @ID create_product
// @Router /product [POST]
// @Summary Create Product
// @Description Create Product
// @Tags Product
// @Accept json
// @Produce json
// @Param product body models.CreateProduct true "CreateProductRequestBody"
// @Success 201 {object} models.Product "GetProductBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) CreateProduct(c *gin.Context) {
	var product models.CreateProduct

	err := c.ShouldBindJSON(&product)
	if err != nil {
		log.Printf("error whiling create: %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Product().Create(context.Background(), &product)
	if err != nil {
		log.Printf("error whiling Create: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling Create").Error())
		return
	}

	resp, err := h.storage.Product().GetByPKey(
		context.Background(),
		&models.ProductPrimarKey{Id: id},
	)

	if err != nil {
		log.Printf("error whiling GetByPKey: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GetByPKey").Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetByIdProduct godoc
// @ID get_by_id_product
// @Router /product/{id} [GET]
// @Summary Get By Id Product
// @Description Get By Id Product
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Product "GetProductBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) GetProductById(c *gin.Context) {

	id := c.Param("id")

	resp, err := h.storage.Product().GetByPKey(
		context.Background(),
		&models.ProductPrimarKey{Id: id},
	)

	if err != nil {
		log.Printf("error whiling GetByPKey: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GetByPKey").Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetListProduct godoc
// @ID get_list_product
// @Router /product [GET]
// @Summary Get List Product
// @Description Get List Product
// @Tags Product
// @Accept json
// @Produce json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} models.GetListProductResponse "GetProductBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) GetProductList(c *gin.Context) {
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

	resp, err := h.storage.Product().GetList(
		context.Background(),
		&models.GetListProductRequest{
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

// UpdateProduct godoc
// @ID update_product
// @Router /product/{id} [PUT]
// @Summary Update Product
// @Description Update Product
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param product body models.UpdateProductSwagger true "CreateProductRequestBody"
// @Success 200 {object} models.Product "GetProductsBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) UpdateProduct(c *gin.Context) {

	var (
		product models.UpdateProduct
	)

	id := c.Param("id")

	if id == "" {
		log.Printf("error whiling update: %v\n", errors.New("required product id").Error())
		c.JSON(http.StatusBadRequest, errors.New("required product id").Error())
		return
	}

	err := c.ShouldBindJSON(&product)
	if err != nil {
		log.Printf("error whiling update: %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	product.Id = id

	rowsAffected, err := h.storage.Product().Update(
		context.Background(),
		&product,
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

	resp, err := h.storage.Product().GetByPKey(
		context.Background(),
		&models.ProductPrimarKey{Id: id},
	)

	if err != nil {
		log.Printf("error whiling GetByPKey: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GetByPKey").Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteByIdProduct godoc
// @ID delete_by_id_product
// @Router /product/{id} [DELETE]
// @Summary Delete By Id Product
// @Description Delete By Id Product
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Product "GetProductBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *HandlerV1) DeleteProduct(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		log.Printf("error whiling update: %v\n", errors.New("required product id").Error())
		c.JSON(http.StatusBadRequest, errors.New("required product id").Error())
		return
	}

	err := h.storage.Product().Delete(
		context.Background(),
		&models.ProductPrimarKey{
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
