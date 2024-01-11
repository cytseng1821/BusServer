package v1

import (
	"BusServer/controllers/tdx"
	"BusServer/postgresql"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

var notifyStopQueue = []string{}
var stopInQueue = map[string]bool{}

type CityBusPlate struct {
	PlateNum     string            `json:"PlateNumb"`
	SubRouteID   string            `json:"SubRouteUID"`
	SubRouteName map[string]string `json:"SubRouteName"`
	Direction    int               `json:"Direction"` // 去返程 : [0:'去程',1:'返程',2:'迴圈',255:'未知']
	StopID       string            `json:"StopUID"`
	StopName     map[string]string `json:"StopName"`
	StopSequence int               `json:"StopSequence"`
	DutyStatus   int               `json:"DutyStatus"`  // 勤務狀態 : [0:'正常',1:'開始',2:'結束']
	BusStatus    int               `json:"BusStatus"`   // 行車狀況 : [0:'正常',1:'車禍',2:'故障',3:'塞車',4:'緊急求援',5:'加油',90:'不明',91:'去回不明',98:'偏移路線',99:'非營運狀態',100:'客滿',101:'包車出租',255:'未知']
	EventType    int               `json:"A2EventType"` // 進站離站 : [0:'離站',1:'進站']
}

func Initialize() {
	updateCityBusNearStop()
	go RoutineUpdateNearStop()
}

func RoutineUpdateNearStop() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[UpdateCityBusNearStop]", err)
		}
	}()

	var ticker = time.NewTicker(1 * time.Minute)
	for range ticker.C {
		updateCityBusNearStop()
	}
}

func updateCityBusNearStop() {
	fmt.Println("Update!!!")
	// 1. 撈出有 user 追蹤的所有 stops 的 routeName
	cityBusRoutes, err := postgresql.GetCityBusRoutesFollowed(context.TODO())
	if err != nil {
		log.Fatalf("updateCityBusNearStop get routes from db: %s", err.Error())
	}
	fmt.Println(cityBusRoutes)

	// 2. for loop 用 routeName 打 TDX api
	ctx := context.TODO()
	var cityBusPlates [][]CityBusPlate
	for _, r := range cityBusRoutes {
		cityBusPointData, statusCode, err := tdx.GetCityBusNearStops(ctx, r.City, r.RouteName)
		if err != nil {
			if statusCode != http.StatusUnauthorized {
				log.Fatalf("updateCityBusNearStop call tdx api: %s", err.Error())
			}
			// refresh token
			if _, err := tdx.GetTDXToken(ctx); err != nil {
				log.Fatalf("updateCityBusNearStop call tdx api: %s", err.Error())
			}
			// call again
			if cityBusPointData, statusCode, err = tdx.GetCityBusNearStops(ctx, r.City, r.RouteName); err != nil {
				log.Fatalf("updateCityBusNearStop call tdx api: %s", err.Error())
			}
		}
		curPlates := make([]CityBusPlate, len(cityBusPointData))
		for i, p := range cityBusPointData {
			curPlates[i] = CityBusPlate{
				PlateNum:     p.PlateNum,
				SubRouteID:   p.SubRouteUID,
				SubRouteName: p.SubRouteName,
				Direction:    p.Direction,
				StopID:       p.StopUID,
				StopName:     p.StopName,
				StopSequence: p.StopSequence,
				DutyStatus:   p.DutyStatus,
				BusStatus:    p.BusStatus,
				EventType:    p.EventType,
			}
			fmt.Printf("\t%s - %s [%s] %v: %s [%s] %v\n", curPlates[i].PlateNum,
				curPlates[i].SubRouteID, curPlates[i].SubRouteName["Zh_tw"], curPlates[i].Direction,
				curPlates[i].StopID, curPlates[i].StopName["Zh_tw"], curPlates[i].StopSequence)
		}
		cityBusPlates = append(cityBusPlates, curPlates)
	}

	// 3. for loop 拿回來的 plate
	for i := range cityBusPlates {
		for _, p := range cityBusPlates[i] {
			// TODO: (1)+(2) cache
			// (1). 用 subRouteID 撈出 stop order by sequence
			cityBusStops, err := postgresql.GetCityBusSubRouteStops(ctx, p.SubRouteID, p.Direction)
			if err != nil {
				log.Fatalf("updateCityBusNearStop get stops from db: %s", err.Error())
			}
			// (2). 建立 stopIDs array 和 map[(subRouteID, stopID)]=sequence
			stopSeq2IDMap, stopID2SeqMap := map[int]string{}, map[string]int{}
			for _, s := range cityBusStops {
				stopSeq2IDMap[s.StopSequence] = s.StopID
				stopID2SeqMap[s.StopID] = s.StopSequence
			}
			fmt.Println(stopSeq2IDMap)
			fmt.Println(stopID2SeqMap)
			// (3). 判斷邏輯和放進 notify event queue 裡面
			curSequence := p.StopSequence
			if curSequence == 0 {
				curSequence = stopID2SeqMap[p.StopID]
			}
			// 當前站檢查是否移除通知
			fmt.Println("cur", curSequence, p.StopID)
			if stopInQueue[p.StopID] {
				stopInQueue[p.StopID] = false
			}
			// 後 5 站檢查是否加入通知
			for j := curSequence + 1; j <= curSequence+5; j++ {
				stopID, ok := stopSeq2IDMap[j]
				if !ok {
					break
				}
				// TODO: 沒有追蹤者的 stop 忽略
				fmt.Println("next", j, stopID)
				if !stopInQueue[stopID] {
					stopInQueue[stopID] = true
					notifyStopQueue = append(notifyStopQueue, stopID)
				}
			}
		}
	}
	fmt.Println(notifyStopQueue)
	fmt.Println(stopInQueue)
}

// subject (stop) & observer (user)

//// api for cronjob per 60 sec
// stopPlates := map[stop]map[plate]bool
// 1. 從 DB 撈取有 user 追蹤的所有 stops (cache)
// 2. 去重複所有 stops 所屬的 routes，並打 TDX api 拿到定點資料
// 3. 掃過定點資料，針對每輛 plate 檢查
//		- 後 5 stops (含當前) 是否有追蹤者
//		i. 沒有追蹤者的 stop 忽略
//		ii. 有追蹤者的 stop
//			if plate 已經在 set 裡面，不做事 -> 因為一輛公車只提醒一次
//			else 把 plate 加進 set 裡面，把 stop 放進 notify event queue 裡面
//		- 前 1 stop 是否有被追蹤
//		i. 沒有追蹤者的 stop 忽略
//		ii. 有被追蹤的 stop
//			if 不在 set 裡面，不做事
//			else 把 plate 從 set 移除

//// backgroud cronjob start from server init
// 當 notify event queue 有資料 (one specific stop) 時
// 1. 從 DB 撈取有追蹤此 stop 的所有 user
// 2. 通知 user（sendgrid 寄信或簡訊）
