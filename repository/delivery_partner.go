package repository

import (
	"errors"

	model "github.com/akshaynanavare/shortest-time/models"
)

type DeliveryPartner interface {
	GetDeliveryPartnerByID(deliveryPartner string) (*model.DeliveryPartner, error)
}

type deliveryPartner struct {
	deliveryPartnerMap map[string]*model.DeliveryPartner
}

func NewDeliveryPartner(DeliveryPartners map[string]*model.DeliveryPartner) DeliveryPartner {
	a := &deliveryPartner{
		deliveryPartnerMap: DeliveryPartners,
	}

	return a
}

func (a *deliveryPartner) GetDeliveryPartnerByID(deliveryPartnerID string) (*model.DeliveryPartner, error) {
	if val, ok := a.deliveryPartnerMap[deliveryPartnerID]; ok {
		return val, nil
	}

	return nil, errors.New("invalid deliveryPartnerID")
}
