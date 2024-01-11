package tdx

import "time"

type TDXBusRouteInfo struct {
	City                string               `json:"City"`
	RouteUID            string               `json:"RouteUID"`
	RouteID             string               `json:"RouteID"`
	RouteName           map[string]string    `json:"RouteName"`
	BusRouteType        int                  `json:"BusRouteType"` // 11:'市區公車',12:'公路客運',13:'國道客運',14:'接駁車'
	DepartureStopName   string               `json:"DepartureStopNameZh"`
	DestinationStopName string               `json:"DestinationStopNameZh"`
	HasSubRoutes        bool                 `json:"HasSubRoutes"`
	SubRoutes           []TDXBusSubRouteInfo `json:"SubRoutes"`
}

type TDXBusSubRouteInfo struct {
	SubRouteUID  string            `json:"SubRouteUID"`
	SubRouteID   string            `json:"SubRouteID"`
	SubRouteName map[string]string `json:"SubRouteName"`
	Direction    int               `json:"Direction"` // 0:'去程',1:'返程',2:'迴圈',255:'未知'
}

type TDXCityBusRoute struct {
	City         string            `json:"City"`
	RouteUID     string            `json:"RouteUID"`
	RouteID      string            `json:"RouteID"`
	RouteName    map[string]string `json:"RouteName"`
	SubRouteUID  string            `json:"SubRouteUID"`
	SubRouteID   string            `json:"SubRouteID"`
	SubRouteName map[string]string `json:"SubRouteName"`
	Direction    int               `json:"Direction"` // 0:'去程',1:'返程',2:'迴圈',255:'未知'
	Stops        []TDXCityBusStop  `json:"Stops"`
}

type TDXCityBusStop struct {
	StopUID      string            `json:"StopUID"`
	StopID       string            `json:"StopID"`
	StopName     map[string]string `json:"StopName"`
	StopSequence int               `json:"StopSequence"`
	StationID    string            `json:"StationID"`
}

type TDXCityBusPointData struct {
	PlateNum      string            `json:"PlateNumb"`
	RouteUID      string            `json:"RouteUID"`
	RouteID       string            `json:"RouteID"`
	RouteName     map[string]string `json:"RouteName"`
	SubRouteUID   string            `json:"SubRouteUID"`
	SubRouteID    string            `json:"SubRouteID"`
	SubRouteName  map[string]string `json:"SubRouteName"`
	Direction     int               `json:"Direction"` // 去返程 : [0:'去程',1:'返程',2:'迴圈',255:'未知']
	StopUID       string            `json:"StopUID"`
	StopID        string            `json:"StopID"`
	StopName      map[string]string `json:"StopName"`
	StopSequence  int               `json:"StopSequence"`
	DutyStatus    int               `json:"DutyStatus"`  // 勤務狀態 : [0:'正常',1:'開始',2:'結束']
	BusStatus     int               `json:"BusStatus"`   // 行車狀況 : [0:'正常',1:'車禍',2:'故障',3:'塞車',4:'緊急求援',5:'加油',90:'不明',91:'去回不明',98:'偏移路線',99:'非營運狀態',100:'客滿',101:'包車出租',255:'未知']
	EventType     int               `json:"A2EventType"` // 進站離站 : [0:'離站',1:'進站']
	SrcUpdateTime time.Time         `json:"SrcUpdateTime"`
	UpdateTime    time.Time         `json:"UpdateTime"`
}

type TDXToken struct {
	AccessToken string `json:"access_token"`
	ExpireIn    int    `json:"expires_in"`
	TokenType   string `json:"Bearer"`
}
