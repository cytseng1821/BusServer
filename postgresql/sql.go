package postgresql

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func InsertCityBusRoutes(ctx context.Context, routes []CityBusRoute) error {
	if len(routes) == 0 {
		return nil
	}

	setMaps := make([]map[string]interface{}, 0, len(routes))
	if err := convertStruct(routes, &setMaps); err != nil {
		return err
	}
	cols := []string{
		"route_id",
		"route_name",
		"subroute_id",
		"subroute_name",
		"direction",
		"city",
		"departure_stop_name",
		"destination_stop_name",
	}
	index := 1
	valueStrs := make([]string, 0, len(routes))
	args := make([]interface{}, 0, len(cols)*len(routes))
	for _, setMap := range setMaps {
		valueStrs = append(valueStrs, "("+generatePlaceHolder(&index, len(cols))+")")
		for _, col := range cols {
			args = append(args, setMap[col])
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO citybus_route
			(%s)
		VALUES
			%s
		ON CONFLICT (route_id, subroute_id, direction)
		DO UPDATE SET
			route_name = EXCLUDED.route_name,
			subroute_name = EXCLUDED.subroute_name,
			city = EXCLUDED.city,
			departure_stop_name = EXCLUDED.departure_stop_name,
			destination_stop_name = EXCLUDED.destination_stop_name`,
		strings.Join(cols, ","),
		strings.Join(valueStrs, ","),
	)
	if _, err := Pool.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func InsertCityBusRoutesStops(ctx context.Context, stops []CityBusStop) error {
	return nil
}

func generatePlaceHolder(index *int, size int) string {
	var placeHolders []string
	for i := 0; i < size; i++ {
		placeHolders = append(placeHolders, fmt.Sprintf("$%d", *index))
		(*index)++
	}
	return strings.Join(placeHolders, ",")
}

func convertStruct(src interface{}, dst interface{}) error {
	jsonStr, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonStr, &dst)
}
