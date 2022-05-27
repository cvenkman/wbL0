package model

import (
	"encoding/json"
	"errors"
)

// get OrderUID from Delivery struct
func GetID(data []byte) (string, error) {
	var model Delivery
	err := json.Unmarshal(data, &model)
	if err != nil {
		return "", errors.New("Unmarshal error: " + err.Error())
	}
	return model.OrderUID, nil
}
