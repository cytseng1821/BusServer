package postgresql

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
)

func InsertCityBusRoutes(ctx context.Context, routes []CityBusRoute, subRoutes []CityBusSubRoute) error {
	if len(routes) == 0 {
		return nil
	}

	// Insert routes first
	setMaps := make([]map[string]interface{}, 0, len(routes))
	if err := convertStruct(routes, &setMaps); err != nil {
		return err
	}
	cols := []string{
		"route_id",
		"route_name",
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
		ON CONFLICT (route_id)
		DO UPDATE SET
			route_name = EXCLUDED.route_name,
			city = EXCLUDED.city,
			departure_stop_name = EXCLUDED.departure_stop_name,
			destination_stop_name = EXCLUDED.destination_stop_name
		RETURNING
			id, route_id`,
		strings.Join(cols, ","),
		strings.Join(valueStrs, ","),
	)
	rows, err := Pool.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var routePKs []struct {
		RoutePK int    `db:"id"`
		RouteID string `db:"route_id"`
	}
	if err := pgxscan.ScanAll(&routePKs, rows); err != nil {
		return err
	}
	routePKMap := make(map[string]int, len(routePKs))
	for _, route := range routePKs {
		routePKMap[route.RouteID] = route.RoutePK
	}

	// Insert subroutes
	setMaps = make([]map[string]interface{}, 0, len(subRoutes))
	if err := convertStruct(subRoutes, &setMaps); err != nil {
		return err
	}
	cols = []string{
		"route_id",
		"subroute_id",
		"subroute_name",
		"direction",
	}
	index = 1
	valueStrs = make([]string, 0, len(subRoutes))
	args = make([]interface{}, 0, len(cols)*len(subRoutes))
	for _, setMap := range setMaps {
		valueStrs = append(valueStrs, "("+generatePlaceHolder(&index, len(cols))+")")
		for _, col := range cols {
			if col == "route_id" {
				args = append(args, routePKMap[setMap[col].(string)])
			} else {
				args = append(args, setMap[col])
			}
		}
	}

	query = fmt.Sprintf(`
		INSERT INTO citybus_subroute
			(%s)
		VALUES
			%s
		ON CONFLICT (subroute_id, direction)
		DO UPDATE SET
			subroute_name = EXCLUDED.subroute_name`,
		strings.Join(cols, ","),
		strings.Join(valueStrs, ","),
	)
	if _, err := Pool.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func InsertCityBusRoutesStops(ctx context.Context, routeID string, stops []CityBusStop, relations []CityBusSubRoute2StopRelation) error {
	if len(stops) == 0 {
		return nil
	}

	// Select subroutes
	query := `
		WITH r AS (
			SELECT id
			FROM citybus_route
			WHERE route_id = $1
		)
		SELECT
			id, subroute_id, direction, route_id
		FROM
			citybus_subroute
		WHERE
			route_id = (SELECT id FROM r)`
	rows, err := Pool.Query(ctx, query, routeID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var subRoutePKs []struct {
		SubRoutePK int    `db:"id"`
		SubRouteID string `db:"subroute_id"`
		Direction  int    `db:"direction"`
		// extend
		RoutePK int `db:"route_id"`
	}
	if err := pgxscan.ScanAll(&subRoutePKs, rows); err != nil {
		return err
	}
	routePK := subRoutePKs[0].RoutePK // all the same route
	subRoutePKMap := make(map[string]int, len(subRoutePKs))
	for _, subRoute := range subRoutePKs {
		subRoutePKMap[subRoute.SubRouteID+"_"+strconv.Itoa(subRoute.Direction)] = subRoute.SubRoutePK
	}

	// Insert stops
	setMaps := make([]map[string]interface{}, 0, len(stops))
	if err := convertStruct(stops, &setMaps); err != nil {
		return err
	}
	cols := []string{
		"route_id",
		"stop_id",
		"stop_name",
	}
	index := 1
	valueStrs := make([]string, 0, len(stops))
	args := make([]interface{}, 0, len(cols)*len(stops))
	for _, setMap := range setMaps {
		valueStrs = append(valueStrs, "("+generatePlaceHolder(&index, len(cols))+")")
		for _, col := range cols {
			if col == "route_id" {
				args = append(args, routePK)
			} else {
				args = append(args, setMap[col])
			}
		}
	}

	query = fmt.Sprintf(`
		INSERT INTO citybus_stop
			(%s)
		VALUES
			%s
		ON CONFLICT (stop_id)
		DO UPDATE SET
			stop_name = EXCLUDED.stop_name
		RETURNING
			id, stop_id`,
		strings.Join(cols, ","),
		strings.Join(valueStrs, ","),
	)
	rows, err = Pool.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var stopPKs []struct {
		StopPK int    `db:"id"`
		StopID string `db:"stop_id"`
	}
	if err := pgxscan.ScanAll(&stopPKs, rows); err != nil {
		return err
	}
	stopPKMap := make(map[string]int, len(stopPKs))
	for _, stop := range stopPKs {
		stopPKMap[stop.StopID] = stop.StopPK
	}

	// Insert Relations
	cols = []string{
		"subroute_id",
		"stop_id",
		"stop_sequence",
	}
	index = 1
	valueStrs = make([]string, 0, len(relations))
	args = make([]interface{}, 0, len(cols)*len(relations))
	for _, relation := range relations {
		valueStrs = append(valueStrs, "("+generatePlaceHolder(&index, len(cols))+")")
		for _, col := range cols {
			if col == "subroute_id" {
				args = append(args, subRoutePKMap[relation.SubRouteID+"_"+strconv.Itoa(relation.Direction)])
			} else if col == "stop_id" {
				args = append(args, stopPKMap[relation.StopID])
			} else {
				args = append(args, relation.StopSequence)
			}
		}
	}

	query = fmt.Sprintf(`
		INSERT INTO citybus_subroute_stop_relation
			(%s)
		VALUES
			%s
		ON CONFLICT (subroute_id, stop_sequence)
		DO UPDATE SET
			stop_id = EXCLUDED.stop_id`,
		strings.Join(cols, ","),
		strings.Join(valueStrs, ","),
	)
	if _, err := Pool.Exec(ctx, query, args...); err != nil {
		return err
	}

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
