package edge_tts

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	winEpoch = 11644473600 // Windows epoch (1601-01-01 00:00:00 UTC)
	sToNs    = 1e9
)

var (
	clockSkewSeconds float64
	skewMutex        sync.Mutex
)

// AdjustClockSkew 调整时钟偏差（线程安全）
func AdjustClockSkew(skewSeconds float64) {
	skewMutex.Lock()
	defer skewMutex.Unlock()
	clockSkewSeconds += skewSeconds
}

// GetUnixTimestamp 获取当前Unix时间戳（含时钟偏差校正）
func GetUnixTimestamp() float64 {
	return float64(time.Now().UTC().UnixNano())/1e9 + clockSkewSeconds
}

// ParseRFC2616Date 解析RFC 2616日期字符串
func ParseRFC2616Date(date string) (float64, error) {
	t, err := time.ParseInLocation(time.RFC1123, date, time.UTC)
	if err != nil {
		return 0, err
	}
	return float64(t.UnixNano()) / 1e9, nil
}

// HandleClientResponseError 处理客户端响应错误
func HandleClientResponseError(resp *http.Response) error {
	serverDate := resp.Header.Get("Date")
	if serverDate == "" {
		return errors.New("no server date in headers")
	}

	serverTime, err := ParseRFC2616Date(serverDate)
	if err != nil {
		return fmt.Errorf("failed to parse server date: %w", err)
	}

	clientTime := GetUnixTimestamp()
	AdjustClockSkew(serverTime - clientTime)
	return nil
}

// GenerateSecMSGec 生成Sec-MS-GEC令牌
func GenerateSecMSGec() string {
	// 获取校正后的时间戳
	ticks := GetUnixTimestamp()

	// 转换为Windows文件时间基准
	ticks += winEpoch

	// 向下取整到最近的5分钟（300秒）
	ticks -= float64(int64(ticks) % 300)

	// 转换为100纳秒间隔
	ticks *= sToNs / 100

	// 生成哈希字符串, 将 ticks 保留9位整数，后面都为0
	strToHash := fmt.Sprintf("%d%s", int(ticks/1e9)*1e9, TRUSTED_CLIENT_TOKEN)
	hash := sha256.Sum256([]byte(strToHash))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}
