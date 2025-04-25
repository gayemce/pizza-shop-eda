package logger

import (
	"os"

	"github.com/gayemce/pizza-shop-eda/order-service/utils"
)

func Log(message any) {
	isLoggedEnabled := os.Getenv("log")
	if isLoggedEnabled != "" {
		utils.AppendToFile("log.txt", message.(string))
	}
}