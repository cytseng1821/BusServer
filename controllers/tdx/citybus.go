package tdx

import (
	"BusServer/constant"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/gin-gonic/gin"
)

var TDXBasicAPI string = "https://tdx.transportdata.tw/api/basic"

func GetCityBusRoutes(c *gin.Context, city, route string) ([]TDXBusRouteInfo, int, error) {
	uri, err := url.Parse(TDXBasicAPI)
	if err != nil {
		return nil, 0, err
	}
	uri.Path = path.Join(uri.Path, "/v2/Bus/Route/City/"+city+"/"+route)

	q := url.Values{}
	q.Add("$format", "JSON")
	uri.RawQuery = q.Encode()

	respBody, statusCode, err := constant.Request(c, constant.RequestParam{
		Method: http.MethodGet,
		URL:    uri.String(),
		Body:   nil,
		Header: http.Header{
			"authorization": []string{"Bearer " + TDXAccessToken},
		},
	})
	if err != nil {
		return nil, statusCode, err
	}

	var respData []TDXBusRouteInfo
	if err = json.Unmarshal(respBody, &respData); err != nil {
		return nil, statusCode, err
	}

	// respBytes, _ := json.MarshalIndent(respData, "", "    ")
	// fmt.Println(string(respBytes))
	return respData, statusCode, nil
}

func GetCityBusRoutesStops(c *gin.Context, city, route string) ([]TDXCityBusRoute, int, error) {
	uri, err := url.Parse(TDXBasicAPI)
	if err != nil {
		return nil, 0, err
	}
	uri.Path = path.Join(uri.Path, "/v2/Bus/StopOfRoute/City/"+city+"/"+route)

	q := url.Values{}
	q.Add("$format", "JSON")
	uri.RawQuery = q.Encode()

	respBody, statusCode, err := constant.Request(c, constant.RequestParam{
		Method: http.MethodGet,
		URL:    uri.String(),
		Body:   nil,
		Header: http.Header{
			"authorization": []string{"Bearer " + TDXAccessToken},
		},
	})
	if err != nil {
		return nil, statusCode, err
	}

	var respData []TDXCityBusRoute
	if err = json.Unmarshal(respBody, &respData); err != nil {
		return nil, statusCode, err
	}

	// respBytes, _ := json.MarshalIndent(respData, "", "    ")
	// fmt.Println(string(respBytes))
	return respData, statusCode, nil
}

func GetCityBusNearStops(c context.Context, city, route string) ([]TDXCityBusPointData, int, error) {
	uri, err := url.Parse(TDXBasicAPI)
	if err != nil {
		return nil, 0, err
	}
	uri.Path = path.Join(uri.Path, "/v2/Bus/RealTimeNearStop/City/"+city+"/"+route)

	q := url.Values{}
	q.Add("$format", "JSON")
	uri.RawQuery = q.Encode()

	respBody, statusCode, err := constant.Request(c, constant.RequestParam{
		Method: http.MethodGet,
		URL:    uri.String(),
		Body:   nil,
		Header: http.Header{
			"authorization": []string{"Bearer " + TDXAccessToken},
		},
	})
	if err != nil {
		return nil, statusCode, err
	}

	var respData []TDXCityBusPointData
	if err = json.Unmarshal(respBody, &respData); err != nil {
		return nil, statusCode, err
	}

	// respBytes, _ := json.MarshalIndent(respData, "", "    ")
	// fmt.Println(string(respBytes))
	return respData, statusCode, nil
}
