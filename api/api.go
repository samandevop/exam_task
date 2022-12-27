package api

import (
	_ "crud/api/docs"
	"crud/api/handler"
	"crud/config"
	"crud/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpApi(cfg *config.Config, r *gin.Engine, storage storage.StorageI) {

	handlerV1 := handler.NewHandlerV1(cfg, storage)

	r.POST("/category", handlerV1.CreateCategory)
	r.GET("/category/:id", handlerV1.GetCategoryById)
	r.GET("/category", handlerV1.GetCategoryList)
	r.PUT("/category/:id", handlerV1.UpdateCategory)
	r.DELETE("/category/:id", handlerV1.DeleteCategory)

	r.POST("/product", handlerV1.CreateProduct)
	r.GET("/product/:id", handlerV1.GetProductById)
	r.GET("/product", handlerV1.GetProductList)
	r.PUT("/product/:id", handlerV1.UpdateProduct)
	r.DELETE("/product/:id", handlerV1.DeleteProduct)

	r.POST("/order", handlerV1.CreateOrder)
	r.GET("/order/:id", handlerV1.GetOrderById)
	r.GET("/order", handlerV1.GetOrderList)
	r.PUT("/order/:id", handlerV1.UpdateOrder)
	r.DELETE("/order/:id", handlerV1.DeleteOrder)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
