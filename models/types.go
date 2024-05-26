package model

import "time"

type Order struct {
	Restaurant *Restaurant
	Customer   *Customer
}

type Restaurant struct {
	ID          string
	Name        string
	AvgPrepTime float64
	Location    *Location
}

type Customer struct {
	ID       string
	Name     string
	Location *Location
}

type Location struct {
	Lat  float64
	Long float64
}

type DeliveryPartner struct {
	ID                         string
	Name                       string
	Order                      []*Order
	TimeRequiredToFinishOrders *time.Time
	CurrentLocation            *Location
	Path                       []string
}
