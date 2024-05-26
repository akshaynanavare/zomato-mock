package utils

import (
	"math"

	model "github.com/akshaynanavare/shortest-time/models"
)

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Haversine function to calculate the distance between two points given their latitude and longitude
func CalculateDistance(p1, p2 *model.Location) float64 {
	var la1, lo1, la2, lo2, r float64

	piRad := math.Pi / 180
	la1 = p1.Lat * piRad
	lo1 = p1.Long * piRad
	la2 = p2.Lat * piRad
	lo2 = p2.Long * piRad

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	meters := 2 * r * math.Asin(math.Sqrt(h))
	kilometers := meters / 1000

	return kilometers
}
