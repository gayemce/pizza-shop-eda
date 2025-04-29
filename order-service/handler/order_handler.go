package handler

import (
	"fmt"

	"github.com/gayemce/pizza-shop-eda/order-service/constants"
	"github.com/gayemce/pizza-shop-eda/order-service/logger"
	"github.com/gayemce/pizza-shop-eda/order-service/service"
	"github.com/gayemce/pizza-shop-eda/order-service/utils"
	"github.com/gin-gonic/gin"
)

var messagePublisherWorker = make(chan map[string]interface{}, 1000)

// 1 worker will process 10 request => 1000

type OrderHandler struct {
	publisher service.IMessagePublisher
}

func (oh *OrderHandler) CreateOrder(ctx *gin.Context) {
	var payload map[string]interface{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Log(fmt.Sprintf("Error mapping body: %v", err))
		ctx.JSON(400, gin.H{
			"message": "Bad Request",
			"statusCode": 400, 
		})
		return
	}
	id := utils.GetId()
	payload["_id"] = id

	messagePublisherWorker <- payload
	
	ctx.JSON(200, gin.H{
		"message": "order is being undertaken, you will message shortly",
		"statusCode": 200,
	})
}

func GetOrderHandler(publisher service.IMessagePublisher) *OrderHandler {
	oh := &OrderHandler{
		publisher: publisher,
	}
	for i := 0; i < cap(messagePublisherWorker); i++ {
		go registerMessagePublisherWorker(i, &publisher)
	}
	return oh
}

func registerMessagePublisherWorker(id int, publisher *service.IMessagePublisher) {
	for message := range messagePublisherWorker {
		err := (*publisher).PublishEvent(constants.TOPIC_ORDER, message)
		if err != nil {
			logger.Log(fmt.Sprintf("Worker %d failed to publish event/message : %v", id, err))
		}
	}
}
