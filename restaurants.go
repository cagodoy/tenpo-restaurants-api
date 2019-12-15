package restaurants

import (
	pb "github.com/cagodoy/tenpo-challenge/lib/proto"
)

// Coord ...
type Coord struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// Restaurant ...
type Restaurant struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Rating         string `json:"rating"`
	Address        string `json:"address"`
	Open           bool   `json:"open"`
	PhotoReference string `json:"photo_reference"`
	Coord          Coord  `json:"coord"`
}

// Service ...
type Service interface {
	ListByCoord(coord Coord, pageToken string) ([]*Restaurant, string, error)
}

// ToProto ...
func (r *Restaurant) ToProto() *pb.Restaurant {
	c := &pb.Coord{
		Latitude:  r.Coord.Latitude,
		Longitude: r.Coord.Latitude,
	}

	return &pb.Restaurant{
		Id:             r.ID,
		Name:           r.Name,
		Rating:         r.Rating,
		Address:        r.Address,
		Open:           r.Open,
		PhotoReference: r.PhotoReference,
		Coord:          c,
	}
}

// FromProto ...
func (r *Restaurant) FromProto(rr *pb.Restaurant) *Restaurant {
	c := Coord{
		Latitude:  rr.GetCoord().GetLatitude(),
		Longitude: rr.GetCoord().GetLongitude(),
	}

	r.ID = rr.Id
	r.Name = rr.Name
	r.Rating = rr.Rating
	r.Address = rr.Address
	r.Open = rr.Open
	r.PhotoReference = rr.PhotoReference
	r.Coord = c

	return r
}
