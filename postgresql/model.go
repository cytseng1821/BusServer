package postgresql

type CityBusRoute struct {
	RouteID             string `json:"route_id"`
	RouteName           string `json:"route_name"`
	City                string `json:"city"`
	DepartureStopName   string `json:"departure_stop_name"`
	DestinationStopName string `json:"destination_stop_name"`
}

type CityBusSubRoute struct {
	SubRouteID   string `json:"subroute_id"`
	SubRouteName string `json:"subroute_name"`
	Direction    int    `json:"direction"`
	// extend
	RouteID string `json:"route_id"`
}

type CityBusStop struct {
	StopID   string `json:"stop_id"`
	StopName string `json:"stop_name"`
	// extend
	RouteID string `json:"route_id"`
}

type CityBusSubRoute2StopRelation struct {
	StopSequence int `json:"stop_sequence"`
	// extend
	SubRouteID string `json:"subroute_id"`
	Direction  int    `json:"direction"`
	StopID     string `json:"stop_id"`
}
