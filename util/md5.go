package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
	"time"
)

var YMRDHS_Format = "2006-01-02 15:04:05"

func Md5String(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return strings.ToLower(hex.EncodeToString(cipherStr))
}

//
func TimeStampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

//
func TimeStampToTimeString(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(YMRDHS_Format)
}

//id
func CreateComplainId(account string) string {
	nowTime := time.Now().UnixNano()
	//  +  md5
	complainId := strings.ToLower(Md5String(account + strconv.FormatInt(nowTime, 10)))
	return complainId
}

func FilePrice(count int64) string {
	var baseCount int64 = 10
	basePrice := "0.001"
	price := "0"
	if count <= 10 {
		price = basePrice
	} else {
		yu := count % baseCount
		priceInt := decimal.NewFromInt(count).Div(decimal.NewFromInt(baseCount)).IntPart()
		if yu > 0 {
			price = decimal.NewFromInt(priceInt).Mul(decimal.RequireFromString(basePrice)).Add(decimal.RequireFromString(basePrice)).String()
		} else {
			price = decimal.NewFromInt(priceInt).Mul(decimal.RequireFromString(basePrice)).String()
		}
	}
	return price
}
