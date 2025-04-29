package routes

import (
	"github.com/gayemce/pizza-shop-eda/order-service/handler"
	"github.com/gayemce/pizza-shop-eda/order-service/service"
	"github.com/gin-gonic/gin"
)

func registerOrderRoutes(r *gin.RouterGroup, publisher service.IMessagePublisher) {

	OrderHandler := handler.GetOrderHandler(publisher)

	r.POST("/create", OrderHandler.CreateOrder,)
}