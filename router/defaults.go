package router

import (
	model "github.com/akshaynanavare/zomato-mock/models"
)

func defaultOrders() map[string][]*model.Order {
	orders := make(map[string][]*model.Order)

	orders["aman"] = []*model.Order{
		{
			Restaurant: &model.Restaurant{
				ID:          "r1",
				Name:        "Star Bazar",
				AvgPrepTime: 10,
				Location: &model.Location{
					Lat:  18.5654,
					Long: 73.8567,
				},
			},
			Customer: &model.Customer{
				ID:   "c1",
				Name: "Akshay",
				Location: &model.Location{
					Lat:  18.4754,
					Long: 73.8567,
				},
			},
		},
		{
			Restaurant: &model.Restaurant{
				ID:          "r2",
				Name:        "Star bucks",
				AvgPrepTime: 30,
				Location: &model.Location{
					Lat:  18.543015,
					Long: 73.786928,
				},
			},
			Customer: &model.Customer{
				ID:   "c2",
				Name: "Akshay 2",
				Location: &model.Location{
					Lat:  18.575026,
					Long: 73.813776,
				},
			},
		},
		{
			Restaurant: &model.Restaurant{
				ID:          "r3",
				Name:        "Sandeep Hotel",
				AvgPrepTime: 8,
				Location: &model.Location{
					Lat:  18.4304,
					Long: 73.8567,
				},
			},
			Customer: &model.Customer{
				ID:   "c3",
				Name: "Akshay 3",
				Location: &model.Location{
					Lat:  18.5204,
					Long: 73.8017,
				},
			},
		},
	}

	return orders
}

func defaultDeliveryPartners() map[string]*model.DeliveryPartner {
	d := map[string]*model.DeliveryPartner{
		"aman": {
			ID:   "a1",
			Name: "Aman",
			CurrentLocation: &model.Location{
				Lat:  18.5204,
				Long: 73.8567,
			},
		},
	}

	return d
}

/*
	Test Data points :
	{18.5204, 73.8567}, // Point A: Central Pune
	{18.5654, 73.8567}, // Point B: 5 km North
	{18.4754, 73.8567}, // Point C: 5 km South
	{18.5204, 73.9117}, // Point D: 5 km East
	{18.5204, 73.8017}, // Point E: 5 km West
	{18.6104, 73.8567}, // Point F: 10 km North
	{18.4304, 73.8567}, // Point G: 10 km South
	{18.5204, 73.9667}, // Point H: 10 km East
	{18.5204, 73.7467}, // Point I: 10 km West

*/
