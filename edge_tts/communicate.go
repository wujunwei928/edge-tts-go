package edge_tts

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Communicate is a struct representing communication with the service.
type Communicate struct {
	text           string
	voice          string
	rate           string
	volume         string
	pitch          string
	proxy          string
	receiveTimeout int
}

type CommunicateOption func(*Communicate) error

// getHeadersAndData returns the headers and data from the given data.
func getHeadersAndData(data interface{}) (map[string][]byte, []byte, error) {
	var dataBytes []byte
	switch v := data.(type) {
	case string:
		dataBytes = []byte(v)
	case []byte:
		dataBytes = v
	default:
		return nil, nil, errors.New("data must be string or []byte")
	}

	headers := make(map[string][]byte)
	headerEnd := bytes.Index(dataBytes, []byte("\r\n\r\n"))
	if headerEnd == -1 {
		return nil, nil, errors.New("invalid data format: no header end")
	}

	headerLines := bytes.Split(dataBytes[:headerEnd], []byte("\r\n"))
	for _, line := range headerLines {
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			return nil, nil, errors.New("invalid header format")
		}
		key := string(bytes.TrimSpace(parts[0]))
		value := bytes.TrimSpace(parts[1])
		headers[key] = value
	}

	return headers, dataBytes[headerEnd+4:], nil
}

// removeIncompatibleCharacters removes incompatible characters from the string.
func removeIncompatibleCharacters(input interface{}) (string, error) {
	var str string
	switch v := input.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return "", errors.New("input must be string or []byte")
	}

	var cleanedStr bytes.Buffer
	for _, char := range str {
		code := int(char)
		if (0 <= code && code <= 8) || (11 <= code && code <= 12) || (14 <= code && code <= 31) {
			cleanedStr.WriteRune(' ')
		} else {
			cleanedStr.WriteRune(char)
		}
	}

	return cleanedStr.String(), nil
}

// connectID generates a UUID without dashes.
func connectID() string {
	u := uuid.New()
	return strings.ReplaceAll(u.String(), "-", "")
}

// splitTextByByteLength splits a string into a list of strings of a given byte length
// while attempting to keep words together.
func splitTextByByteLength(text []byte, byteLength int) [][]byte {
	var result [][]byte

	if byteLength <= 0 {
		panic("byteLength must be greater than 0")
	}

	for len(text) > byteLength {
		// Find the last space in the string within the byte length
		splitAt := byteLength
		for i := byteLength; i > 0; i-- {
			if utf8.RuneStart(text[i]) {
				splitAt = i
				break
			}
		}

		// Verify all & are terminated with a ;
		for bytes.Contains(text[:splitAt], []byte("&")) {
			ampersandIndex := bytes.LastIndex(text[:splitAt], []byte("&"))
			if bytes.Index(text[ampersandIndex:splitAt], []byte(";")) != -1 {
				break
			}

			splitAt = ampersandIndex - 1
			if splitAt < 0 {
				panic("Maximum byte length is too small or invalid text")
			}
			if splitAt == 0 {
				break
			}
		}

		// Append the string to the list
		newText := bytes.TrimSpace(text[:splitAt])
		if len(newText) > 0 {
			result = append(result, newText)
		}

		text = text[splitAt:]
	}

	newText := bytes.TrimSpace(text)
	if len(newText) > 0 {
		result = append(result, newText)
	}

	return result
}

// mkSSML creates an SSML string from the given parameters.
func mkSSML(text string, voice string, rate string, volume string, pitch string) string {
	ssml := fmt.Sprintf(
		"<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='en-US'>"+
			"<voice name='%s'><prosody pitch='%s' rate='%s' volume='%s'>%s</prosody></voice></speak>", voice, pitch, rate, volume, text)
	return ssml
}

// dateToString returns a JavaScript-style date string.
func dateToString() string {
	utcTime := time.Now().UTC()
	timeString := utcTime.Format("Mon Jan 02 2006 15:04:05 GMT+0000 (Coordinated Universal Time)")
	return timeString
}

// ssmlHeadersPlusData returns the headers and data to be used in the request.
func ssmlHeadersPlusData(requestID string, timestamp string, ssml string) string {
	headersAndData := fmt.Sprintf(
		"X-RequestId:%s\r\n"+
			"Content-Type:application/ssml+xml\r\n"+
			"X-Timestamp:%sZ\r\n"+
			"Path:ssml\r\n\r\n"+
			"%s",
		requestID, timestamp, ssml)
	return headersAndData
}

// calcMaxMesgSize calculates the maximum message size for the given voice, rate, and volume.
func calcMaxMesgSize(voice string, rate string, volume string, pitch string) int {
	websocketMaxSize := 1 << 16
	// Calculate overhead per message
	overheadPerMessage := len(ssmlHeadersPlusData(connectID(), dateToString(), mkSSML("", voice, rate, volume, pitch))) + 50 // margin of error
	return websocketMaxSize - overheadPerMessage
}

// ValidateStringParam validates the given string parameter based on type and pattern.
func ValidateStringParam(paramName, paramValue, pattern string) (string, error) {
	if len(paramValue) == 0 {
		return "", errors.New(fmt.Sprintf("%s不能为空", paramName))
	}
	match, err := regexp.MatchString(pattern, paramValue)
	if err != nil {
		return "", err
	}
	if !match {
		return "", errors.New(fmt.Sprintf("%s不符合模式%s", paramName, pattern))
	}
	return paramValue, nil
}

// NewCommunicate initializes the Communicate struct.
func NewCommunicate(text string, options ...CommunicateOption) (*Communicate, error) {
	c := &Communicate{
		text:           text,
		voice:          "Microsoft Server Speech Text to Speech Voice (en-US, AriaNeural)",
		rate:           "+0%",
		volume:         "+0%",
		pitch:          "+0Hz",
		proxy:          "",
		receiveTimeout: 10,
	}

	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// SetVoice sets the voice for communication.
func SetVoice(voice string) CommunicateOption {
	return func(c *Communicate) error {
		//var err error
		//c.voice, err = ValidateStringParam("voice", voice, `^Microsoft Server Speech Text to Speech Voice \(.+,.+\)$`)
		//return err
		c.voice = voice
		return nil
	}
}

// SetRate sets the rate for communication.
func SetRate(rate string) CommunicateOption {
	return func(c *Communicate) error {
		var err error
		c.rate, err = ValidateStringParam("rate", rate, `^[+-]\d+%$`)
		return err
	}
}

// SetVolume sets the volume for communication.
func SetVolume(volume string) CommunicateOption {
	return func(c *Communicate) error {
		var err error
		c.volume, err = ValidateStringParam("volume", volume, `^[+-]\d+%$`)
		return err
	}
}

// SetPitch sets the pitch for communication.
func SetPitch(pitch string) CommunicateOption {
	return func(c *Communicate) error {
		var err error
		c.pitch, err = ValidateStringParam("pitch", pitch, `^[+-]\d+(Hz|%)$`)
		return err
	}
}

// SetProxy sets the proxy for communication.
func SetProxy(proxy string) CommunicateOption {
	return func(c *Communicate) error {
		c.proxy = proxy
		return nil
	}
}

// SetReceiveTimeout sets the receive timeout for communication.
func SetReceiveTimeout(receiveTimeout int) CommunicateOption {
	return func(c *Communicate) error {
		c.receiveTimeout = receiveTimeout
		return nil
	}
}

func (c *Communicate) newWebSocketConn() (*websocket.Conn, error) {
	dialer := &websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: true,
	}

	if len(c.proxy) > 0 {
		proxyURL, err := url.Parse(c.proxy)
		if err != nil {
			return nil, err
		}
		dialer.Proxy = http.ProxyURL(proxyURL)
	}

	header := http.Header{}
	for k, v := range WSS_HEADERS {
		header.Set(k, v)
	}

	dialCtx, dialContextCancel := context.WithTimeout(context.Background(), time.Duration(c.receiveTimeout)*time.Second)
	defer func() {
		dialContextCancel()
	}()

	reqUrl := fmt.Sprintf("%s&Sec-MS-GEC=%s&Sec-MS-GEC-Version=%s&ConnectionId=%s", WSS_URL, GenerateSecMSGec(), SEC_MS_GEC_VERSION, connectID())
	conn, _, err := dialer.DialContext(dialCtx, reqUrl, header)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *Communicate) Stream() ([]byte, error) {
	conn, err := c.newWebSocketConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var finished = make(chan struct{})
	var failed = make(chan error)
	audioData := make([]byte, 0)
	go func() {
		defer func() {
			close(finished)
			close(failed)
		}()
		for {
			receivedType, receivedData, receivedErr := conn.ReadMessage()
			if receivedType == -1 && receivedData == nil && receivedErr != nil { //已经断开链接
				failed <- receivedErr
				return
			}

			switch receivedType {
			case websocket.TextMessage:
				textHeader, _, textErr := getHeadersAndData(receivedData)
				if textErr != nil {
					failed <- textErr
					return
				}
				if string(textHeader["Path"]) == "turn.end" {
					return
				}
			case websocket.BinaryMessage:
				if len(receivedData) < 2 {
					failed <- errors.New("we received a binary message, but it is missing the header length")
					return
				}

				headerLength := binary.BigEndian.Uint16(receivedData[:2])
				if len(receivedData) < int(headerLength+2) {
					failed <- errors.New("we received a binary message, but it is missing the audio data")
					return
				}

				audioData = append(audioData, receivedData[2+headerLength:]...)
			default:
				log.Println("recv:", receivedData)
			}
		}
	}()

	err = conn.WriteMessage(websocket.TextMessage, []byte(c.getCommandRequestContent()))
	if err != nil {
		return nil, err
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(c.getSSMLRequestContent()))
	if err != nil {
		return nil, err
	}

	select {
	case <-finished:
		return audioData, err
	case errMessage := <-failed:
		return nil, errMessage
	}
}

func (c *Communicate) getCommandRequestContent() string {
	var builder strings.Builder

	// 拼接X-Timestamp部分
	builder.WriteString(fmt.Sprintf("X-Timestamp:%s\r\n", dateToString()))

	// 拼接Content-Type部分
	builder.WriteString("Content-Type:application/json; charset=utf-8\r\n")

	// 拼接Path部分
	builder.WriteString("Path:speech.config\r\n\r\n")

	// 拼接JSON部分
	builder.WriteString(`{"context":{"synthesis":{"audio":{"metadataoptions":{`)
	builder.WriteString(`"sentenceBoundaryEnabled":"false","wordBoundaryEnabled":"true"},`)
	builder.WriteString(`"outputFormat":"audio-24khz-48kbitrate-mono-mp3"`)
	builder.WriteString("}}}}\r\n")

	return builder.String()
}

func (c *Communicate) getSSMLRequestContent() string {
	return ssmlHeadersPlusData(
		connectID(),
		dateToString(),
		mkSSML(c.text, c.voice, c.rate, c.volume, c.pitch),
	)
}
