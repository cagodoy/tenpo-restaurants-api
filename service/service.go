package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	restaurants "github.com/cagodoy/tenpo-restaurants-api"
	"googlemaps.github.io/maps"

	history "github.com/cagodoy/tenpo-history-api"
	nats "github.com/nats-io/nats.go"
)

// NewRestaurants ...
func NewRestaurants(conn *nats.EncodedConn) *Restaurants {
	return &Restaurants{
		Nats: conn,
	}
}

// Restaurants ...
type Restaurants struct {
	Nats *nats.EncodedConn
}

// ListByCoord ...
func (us *Restaurants) ListByCoord(coord restaurants.Coord, pageToken string) ([]*restaurants.Restaurant, string, error) {
	// get API_KEY env value
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("missing env variable API_KEY, using default value...")
		os.Exit(1)
	}

	// create new google maps client
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, "", err
	}

	// parse cood values
	loc, err := maps.ParseLatLng(coord.Latitude + "," + coord.Longitude)
	if err != nil {
		return nil, "", err
	}

	// prepare request
	r := &maps.NearbySearchRequest{
		// TODO(ca): get radius value from req param
		Radius:   5000,
		Keyword:  "restaurant",
		Language: "spanish",
		Location: &loc,
	}

	// add pageToken value if present
	if pageToken != "" {
		r.PageToken = pageToken
	}

	// fetch request to google maps api
	resp, err := client.NearbySearch(context.Background(), r)

	// prepare result in slice
	rr := make([]*restaurants.Restaurant, 0)

	// parse response
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

	// emit event to history service
	go func() {
		he := &history.CreateHistoryEvent{
			//TODO(ca): get user_id in req param
			UserID:    "123456",
			Latitude:  coord.Latitude,
			Longitude: coord.Latitude,
		}

		us.Nats.Publish("history.create", he)
		log.Println("Published to History.create service")
	}()

	return rr, resp.NextPageToken, nil
}
