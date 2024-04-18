package edge_tts

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Voice struct {
	Name           string   `json:"Name"`
	ShortName      string   `json:"ShortName"`
	Gender         string   `json:"Gender"`
	Locale         string   `json:"Locale"`
	SuggestedCodec string   `json:"SuggestedCodec"`
	FriendlyName   string   `json:"FriendlyName"`
	Status         string   `json:"Status"`
	VoiceTag       VoiceTag `json:"VoiceTag"`
}

type VoiceTag struct {
	ContentCategories  []string `json:"ContentCategories"`
	VoicePersonalities []string `json:"VoicePersonalities"`
}

func ListVoices(proxyURL string) ([]Voice, error) {
	client := resty.New()
	if len(proxyURL) > 0 {
		client.SetProxy(proxyURL)
	}

	resp, err := client.R().
		EnableTrace().
		SetHeaders(map[string]string{
			"Authority":        "speech.platform.bing.com",
			"Sec-CH-UA":        `" Not;A Brand";v="99", "Microsoft Edge";v="91", "Chromium";v="91"`,
			"Sec-CH-UA-Mobile": "?0",
			"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36 Edg/91.0.864.41",
			"Accept":           "*/*",
			"Sec-Fetch-Site":   "none",
			"Sec-Fetch-Mode":   "cors",
			"Sec-Fetch-Dest":   "empty",
			"Accept-Encoding":  "gzip, deflate, br",
			"Accept-Language":  "en-US,en;q=0.9",
		}).
		Get(VOICE_LIST_URL)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list voices, http status code: %s", resp.Status())
	}

	var voices []Voice
	err = json.Unmarshal(resp.Body(), &voices)
	if err != nil {
		return nil, err
	}

	return voices, nil
}
