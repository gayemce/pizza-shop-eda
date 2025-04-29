package routes

import (
	"github.com/gayemce/pizza-shop-eda/order-service/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, publisher service.IMessagePublisher) {
	orderRoutes := r.Group("/order-service")
	{
		registerOrderRoutes(orderRoutes, publisher)
	}
}