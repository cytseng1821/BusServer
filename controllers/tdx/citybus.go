package tdx

import (
	"BusServer/constant"
	"encoding/json"
	"fmt"
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

func GetCityBusNearStops(c *gin.Context, city, route string) {
	uri, err := url.Parse(TDXBasicAPI)
	if err != nil {
		fmt.Println("parse", err.Error())
		return
	}
	uri.Path = path.Join(uri.Path, "/v2/Bus/RealTimeNearStop/City/"+city+"/"+route)

	q := url.Values{}
	q.Add("$format", "JSON")
	uri.RawQuery = q.Encode()
	fmt.Println(uri.String())

	respBody, _, err := constant.Request(c, constant.RequestParam{
		Method: http.MethodGet,
		URL:    uri.String(),
		Body:   nil,
		Header: http.Header{
			"authorization": []string{"Bearer " + TDXAccessToken},
		},
	})
	if err != nil {
		fmt.Println("request", err.Error())
		return
	}

	var respData []TDXCityBusPointData
	if err = json.Unmarshal(respBody, &respData); err != nil {
		fmt.Println("unmarshal", err.Error(), string(respBody))
		return
	}

	respBytes, _ := json.MarshalIndent(respData, "", "    ")
	fmt.Println(string(respBytes))
}
