# Restaurants-API

Microservice implemented in Golang that get nearby restaurants from Google Places API.

## GRPC Service

```go
service RestaurantService {
  rpc ListByCoord(RestaurantListByCoordRequest) returns (RestaurantListByCoordResponse) {}
}

message Restaurant {
  string id = 1;
  string name = 2;
  string rating = 3;
  string address = 4;
  bool open = 5;
  string photo_reference = 6;
  Coord coord = 7;
}

message RestaurantListByCoordRequest {
  Coord coord = 1;
  string page_token = 2;
}

message MetaRestaurantListByCoord {
  string page_token = 1;
}

message RestaurantListByCoordResponse {
  repeated Restaurant data = 1;
  MetaRestaurantListByCoord meta = 2;
  Error error = 3;
}
```

## Commands (Development)

`make build`: build restaurants service for osx.

`make linux`: build restaurants service for linux os.

`make docker .`: build docker.

`docker run -it -p 5030:5030 tenpo-restaurants-api`: run docker.

`PORT=<port> API_KEY=<api_key> make r`: run tenpo restaurants service.
