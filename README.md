# CouponRushSystem
===

# 題目

設計一套搶購機制，規格：

1. 每天定點23點發送優惠券，需要用戶提前1-5分鐘前先預約，優惠券的數量為預約用戶數量的20%，
2. 搶購優惠券時，每個用戶只能搶一次，只有1分鐘的搶購時間，怎麼盡量保證用戶搶到的概率盡量貼近20%。
3. 需要考慮有300人搶和30000人搶的場景。
4. 設計表結構和索引、編寫主要程式碼。

# 交付時限/方式

- [ ]  題目送出後 7 天內
- [ ]  交付方式不限，可以是 Github, 也可以是 PDF.


# 設計

## 預約發放優惠卷 

每日由後台排程新增 `active_YYYYMMDD` 的 key 進入 redis. 執行 `SET active_20240424 0` .

開放時間為每日 22:55 - 22:59 , 給用戶進行預約活動的登記; 其他時間一率回應 `403` Forbidden.

設計一個bloom filter, 將user id 放入, 不存在bloom filter的用戶才能預約 (反之, 已預約的不會不存在)
換個思路 用redis bitmaps, 將user id 對應的bit設置為1. 這樣有沒有預約過就會知道. key採用今天日期.  user id太長超過2^32次方在來煩惱.

再來使用取號機, 每個預約成功請求給個 `預約號碼牌`數量 ; 因為優惠卷數量是預約人數的 20%, 所以每取號5張(或每次取號牌尾數是0) 就即時產生一組優惠卷號碼於 rdb/redis.
`INCR active_YYYYMMDD`


## 搶購優惠卷

已預約成功的用戶, 可以在 23:00 - 23:01 進行搶購優惠卷的行為; i.e. 一分鐘內的秒殺搶購流量大量湧入.

用戶搶購請求進來, 首先要驗證是否是已預約成功的用戶.(思考, 有沒有可能不用去資料庫驗證, 而靠簽章驗證是已預約成功; 有困難, 搞不好用戶登出過或換裝置);
但bloom filter 只有對不存在絕對正確, 找存在是偽陽性. 
改用 user id 查找redis bitmaps.

已預約成功的用戶, 按照請求順序依序購買到優惠卷; 換句話說比較慢進來搶購的用戶可能就買不到了. (避免超賣問題).

當有第一個用戶請求發現已經沒優惠卷了, 此時 應用程式內的 flag 可以改成 `售完`. 以快速回應客戶搶購請求, 也避免給資料庫壓力.


購買優惠卷 solutions:

1. 於Postgres中 , 因為能請求的已經是有預約成功的用戶了, 可以篩選掉大部分無意義的搶購請求.
```
tx.Begin()

UPDATE coupon SET buyer=$user_id , status=reserved WHERE status=unused LIMIT 1 returning id;

if id > 0 {
    // insert a record to reserved_coupon  table with user id and coupon id
} else {
    // update flag state in memory
}

tx.Commit()
```

2. 於Redis中

coupon預熱加載到redis的list中. 

寫lua script, 有預約請求, 就把coupon list給pop, 然後塞進去另一個 coupon_reserved_list中, 並且也對 coupon_YYYYMMDD_reserved 針對user id 的bit設定成1

當結果取得coupon或者取得`ALREADY_RECEIVED`時, 該服務也能keep 一份user coupon cache with TTL 2分鐘.  同一個用戶的重複請求就能在application 直接reject.


```lua
-- Lua script for handling coupon distribution in Redis
-- KEYS[1] = list of available coupons (e.g., "coupons_20240424")
-- KEYS[2] = list of reserved coupons (e.g., "coupons_20240424_reserved")
-- KEYS[3] = bitmap for tracking users who have received a coupon (e.g., "coupon_users_20240424")
-- ARGV[1] = user ID

-- Check if the user has already got a coupon
if redis.call('GETBIT', KEYS[3], ARGV[1]) == 1 then
    return 'ALREADY_RECEIVED' -- User has already received a coupon, return specific message
end

-- Try to pop a coupon from the available list
local coupon = redis.call('LPOP', KEYS[1])
if coupon then
    -- Mark the user as having received a coupon
    redis.call('SETBIT', KEYS[3], ARGV[1], 1)
    -- Push the coupon to the reserved list
    redis.call('RPUSH', KEYS[2], coupon)
    -- Return the coupon code to the user
    return coupon
else
    -- No coupons left, return specific message
    return 'NO_COUPONS'
end
```

送個event 到mq, 慢慢同步至postgres,
event record, coupon id + user id + event time + event type 'received_coupon`

再有一些consumer, prefetch能多點, 來同步這些event 到資料庫並更新優惠卷狀態

## HTTP Status

403 Forbidden


## Busniess Error Codes

NOT_IN_ACTIVE_TIME

ERROR_DUPLICATED_RESERVATION 

NOT_RESERVATION


## Data Schema


# Project Layout

## /cmd

執行 Server 或 Preload 的進入點

## /doc

存放 openapi doc 

## /config

存放啟動該
