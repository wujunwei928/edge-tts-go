package edge_tts

import (
	"fmt"
	"strings"
)

const (
	PackageVersion = "0.0.1" // edge-tts-go 包版本
)

// edge tts 相关接口
const (
	TRUSTED_CLIENT_TOKEN = "6A5AA1D4EAFF4E9FB37E23D68491D6F4"
	WSS_URL              = "wss://speech.platform.bing.com/consumer/speech/synthesize/readaloud/edge/v1?TrustedClientToken=" + TRUSTED_CLIENT_TOKEN
	VOICE_LIST_URL       = "https://speech.platform.bing.com/consumer/speech/synthesize/readaloud/voices/list?trustedclienttoken=" + TRUSTED_CLIENT_TOKEN
)

const CHROMIUM_FULL_VERSION = "134.0.3124.66"

var (
	CHROMIUM_MAJOR_VERSION = strings.SplitN(CHROMIUM_FULL_VERSION, ".", 2)[0]
	SEC_MS_GEC_VERSION     = fmt.Sprintf("1-%s", CHROMIUM_FULL_VERSION)

	BASE_HEADERS = map[string]string{
		"User-Agent":      fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s.0.0.0 Safari/537.36 Edg/%s.0.0.0", CHROMIUM_MAJOR_VERSION, CHROMIUM_MAJOR_VERSION),
		"Accept-Encoding": "gzip, deflate, br",
		"Accept-Language": "en-US,en;q=0.9",
	}

	WSS_HEADERS = MergeMap(map[string]string{
		"Pragma":        "no-cache",
		"Cache-Control": "no-cache",
		"Origin":        "chrome-extension://jdiccldimpdaibmpdkjnbmckianbfold",
	}, BASE_HEADERS)

	VOICE_HEADERS = MergeMap(map[string]string{
		"Authority":        "speech.platform.bing.com",
		"Sec-CH-UA":        fmt.Sprintf(`" Not;A Brand";v="99", "Microsoft Edge";v="%s", "Chromium";v="%s"`, CHROMIUM_MAJOR_VERSION, CHROMIUM_MAJOR_VERSION),
		"Sec-CH-UA-Mobile": "?0",
		"Accept":           "*/*",
		"Sec-Fetch-Site":   "none",
		"Sec-Fetch-Mode":   "cors",
		"Sec-Fetch-Dest":   "empty",
	}, BASE_HEADERS)
)

func MergeMap(m1, m2 map[string]string) map[string]string {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}
