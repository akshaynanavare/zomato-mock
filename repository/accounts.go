package repository

import (
	"errors"

	model "github.com/akshaynanavare/shortest-time/models"
)

type Order interface {
	GetOrdersByID(deliveryPartner string) ([]*model.Order, error)
}

type order struct {
	orderMap map[string][]*model.Order
}

func NewOrder(orders map[string][]*model.Order) Order {
	a := &order{
		orderMap: orders,
	}

	return a
}

func (a *order) GetOrdersByID(deliveryPartnerID string) ([]*model.Order, error) {
	if val, ok := a.orderMap[deliveryPartnerID]; ok {
		return val, nil
	}

	return nil, errors.New("invalid deliveryPartnerID")
}
