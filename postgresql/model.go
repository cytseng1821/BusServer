package postgresql

type CityBusRoute struct {
	RouteID             string `json:"route_id"`
	RouteName           string `json:"route_name"`
	SubRouteID          string `json:"subroute_id"`
	SubRouteName        string `json:"subroute_name"`
	Direction           int    `json:"direction"`
	City                string `json:"city"`
	DepartureStopName   string `json:"departure_stop_name"`
	DestinationStopName string `json:"destination_stop_name"`
}

type CityBusStop struct {
	RoutePK      int    `json:"route_pk"`
	StopID       string `json:"stop_id"`
	StopName     string `json:"stop_name"`
	StopSequence int    `json:"stop_sequence"`
	// extend
	RouteID    string `json:"route_id"`
	SubRouteID string `json:"subroute_id"`
	Direction  int    `json:"direction"`
}
