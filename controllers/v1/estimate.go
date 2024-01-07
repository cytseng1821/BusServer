package v1

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
