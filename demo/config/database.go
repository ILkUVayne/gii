package config

import (
	"gii/demo/model"
)

var modelMaps = map[string]interface{}{
	"user":  &model.User{},
	"order": &model.Order{},
}

func CheckTable() {
	for _, m := range modelMaps {
		if m1, ok := m.(model.ICheckTable); ok {
			m1.CheckTableExist()
		}
	}
}
