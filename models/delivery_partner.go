package model

import (
	"time"
)

type DeliveryPartner struct {
	ID                         string
	Name                       string
	Order                      []*Order
	TimeRequiredToFinishOrders *time.Time
	CurrentLocation            *Location
	Path                       []string
}
