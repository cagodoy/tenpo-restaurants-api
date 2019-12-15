package restaurantssvc

import (
	"context"
	"fmt"

	pb "github.com/cagodoy/tenpo-challenge/lib/proto"
	restaurants "github.com/cagodoy/tenpo-restaurants-api"
	"github.com/cagodoy/tenpo-restaurants-api/service"

	"googlemaps.github.io/maps"
)

var _ pb.RestaurantServiceServer = (*Service)(nil)

// Service ...
type Service struct {
	restaurantsSvc restaurants.Service
}

// New ...
func New() *Service {
	return &Service{
		restaurantsSvc: service.NewRestaurants(),
	}
}

// ListByCoord List nearby restaurants by coord.
func (as *Service) ListByCoord(ctx context.Context, gr *pb.RestaurantListByCoordRequest) (*pb.RestaurantListByCoordResponse, error) {
	lat := gr.GetCoord().GetLatitude()
	lng := gr.GetCoord().GetLongitude()

	_, err := maps.ParseLatLng(lat + "," + lng)
	if err != nil {
		return &pb.RestaurantListByCoordResponse{
			Data: nil,
			Meta: nil,
			Error: &pb.Error{
				Code:    500,
				Message: "invalid coord values",
			},
		}, nil
	}

	c := restaurants.Coord{
		Latitude:  lat,
		Longitude: lng,
	}

	pageToken := gr.GetPageToken()

	listedRestaurants, nextPageToken, err := as.restaurantsSvc.ListByCoord(c, pageToken)
	if err != nil {
		fmt.Println(fmt.Sprintf("[GRPC][RestaurantsService][ListByCoord][Error] %v", err))
		return &pb.RestaurantListByCoordResponse{
			Data: nil,
			Meta: nil,
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
		Data: data,
		Meta: &pb.MetaRestaurantListByCoord{
			PageToken: nextPageToken,
		},
		Error: nil,
	}

	fmt.Println(fmt.Sprintf("[GRPC][RestaurantsService][List][Response] Listed %v restaurants", len(res.Data)))
	return res, nil
}
