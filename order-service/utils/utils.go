package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AppendToFile(filePath string, data interface{}) error {

	var dataStr string

	switch v := data.(type) {
	case string:
		dataStr = v
	default:
		jsonData, err := json.Marshal(data)
		if err != nil {
			dataStr = fmt.Sprintf("%v", data)
		} else {
			dataStr = string(jsonData)
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(dataStr)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}
	return nil
}

func GetId() interface{} {
	return primitive.NewObjectID()
}

func GetIdString(objId interface{}) string {
	id := objId.(primitive.ObjectID)
	if id == primitive.NilObjectID {
		return ""
	}
	return id.Hex()
}