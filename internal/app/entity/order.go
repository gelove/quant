package entity

import (
	"quant/internal/app/entity/common"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	common.Order
}
