package restaurantssvc

import (
	"context"
	"fmt"
	"log"

	pb "github.com/cagodoy/tenpo-challenge/lib/proto"
	restaurants "github.com/cagodoy/tenpo-restaurants-api"
	"github.com/cagodoy/tenpo-restaurants-api/service"

	"googlemaps.github.io/maps"

	nats "github.com/nats-io/nats.go"
)

var _ pb.RestaurantServiceServer = (*Service)(nil)

// Service ...
type Service struct {
	restaurantsSvc restaurants.Service
}

// New ...
func New(conn *nats.EncodedConn) *Service {
	return &Service{
		restaurantsSvc: service.NewRestaurants(conn),
	}
}

// ListByCoord List nearby restaurants by coord.
func (as *Service) ListByCoord(ctx context.Context, gr *pb.RestaurantListByCoordRequest) (*pb.RestaurantListByCoordResponse, error) {
	lat := gr.GetCoord().GetLatitude()
	lng := gr.GetCoord().GetLongitude()

	c := restaurants.Coord{
		Latitude:  lat,
		Longitude: lng,
	}

	_, err := maps.ParseLatLng(c.GetLatLngStr())
	if err != nil {
		return &pb.RestaurantListByCoordResponse{
			Data: nil,
			Error: &pb.Error{
				Code:    500,
				Message: "invalid coord values",
			},
		}, nil
	}

	userID := gr.GetUserId()

	listedRestaurants, err := as.restaurantsSvc.ListByCoord(c, userID)
	if err != nil {
		log.Println(fmt.Sprintf("[GRPC][RestaurantsService][ListByCoord][Error] %v", err))
		return &pb.RestaurantListByCoordResponse{
			Data: nil,
			Error: &pb.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	data := make([]*pb.Restaurant, 0)
	for _, restaurant := range listedRestaurants {
		data = append(data, restaurant.ToProto())
	}

	res := &pb.RestaurantListByCoordResponse{
		Data:  data,
		Error: nil,
	}

	log.Println(fmt.Sprintf("[GRPC][RestaurantsService][List][Response] Listed %v restaurants", len(res.Data)))
	return res, nil
}
