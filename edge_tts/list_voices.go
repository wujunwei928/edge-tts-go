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

	voiceListUrl := fmt.Sprintf(
		"%s&Sec-MS-GEC=%s&Sec-MS-GEC-Version=%s",
		VOICE_LIST_URL,
		GenerateSecMSGec(),
		SEC_MS_GEC_VERSION,
	)

	resp, err := client.R().
		EnableTrace().
		SetHeaders(VOICE_HEADERS).
		Get(voiceListUrl)

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
