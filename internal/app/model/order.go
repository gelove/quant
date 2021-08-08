package model

import (
	"quant/internal/app/entity"
	"quant/pkg/orm"

	"github.com/pkg/errors"
)

type Order struct{}

func (o *Order) Get(symbol string) ([]*entity.Order, error) {
	orders := make([]*entity.Order, 0, 1<<2)
	err := orm.DB.Where("symbol = ?", symbol).Where("status = ?", "NEW").Find(&orders).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return orders, nil
}

func (o *Order) Find(symbol string, orderID int64) (*entity.Order, error) {
	order := &entity.Order{}
	err := orm.DB.Where("symbol = ?", symbol).Where("orderId = ?", orderID).First(order).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return order, nil
}

func (o *Order) Create(order *entity.Order) error {
	err := orm.DB.Create(order).Error
	return errors.WithStack(err)
}

func (o *Order) Save(order *entity.Order) error {
	err := orm.DB.Save(order).Error
	return errors.WithStack(err)
}
