package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	restaurants "github.com/cagodoy/tenpo-restaurants-api"
	"googlemaps.github.io/maps"
)

// NewRestaurants ...
func NewRestaurants() *Restaurants {
	return &Restaurants{}
}

// Restaurants ...
type Restaurants struct{}

// ListByCoord ...
func (us *Restaurants) ListByCoord(coord restaurants.Coord, pageToken string) ([]*restaurants.Restaurant, string, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("missing env variable API_KEY, using default value...")
		os.Exit(1)
	}

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, "", err
	}

	loc, err := maps.ParseLatLng(coord.Latitude + "," + coord.Longitude)
	if err != nil {
		return nil, "", err
	}

	r := &maps.NearbySearchRequest{
		Radius:   5000,
		Keyword:  "restaurant",
		Language: "spanish",
		Location: &loc,
	}

	if pageToken != "" {
		r.PageToken = pageToken
	}

	resp, err := client.NearbySearch(context.Background(), r)

	rr := make([]*restaurants.Restaurant, 0)

	for _, res := range resp.Results {
		c := restaurants.Coord{
			Latitude:  strconv.FormatFloat(res.Geometry.Location.Lat, 'f', -1, 64),
			Longitude: strconv.FormatFloat(res.Geometry.Location.Lng, 'f', -1, 64),
		}

		restaurant := &restaurants.Restaurant{
			ID:             res.PlaceID,
			Name:           res.Name,
			Address:        res.Vicinity,
			Rating:         fmt.Sprintf("%f", res.Rating),
			Open:           *res.OpeningHours.OpenNow,
			PhotoReference: res.Photos[0].PhotoReference,
			Coord:          c,
		}

		rr = append(rr, restaurant)
	}

	return rr, resp.NextPageToken, nil
}
