package edge_tts

import (
	"fmt"
)

// 预定义输出格式常量
const (
	OutputFormatMP3   = "mp3"    // MP3 48kbps（默认）
	OutputFormatMP3HQ = "mp3-hq" // MP3 96kbps
	OutputFormatWebM  = "webm"   // WebM Opus
)

// edgeOutputFormats 合法格式到 Edge API 原始格式字符串的映射
var edgeOutputFormats = map[string]string{
	OutputFormatMP3:   "audio-24khz-48kbitrate-mono-mp3",
	OutputFormatMP3HQ: "audio-24khz-96kbitrate-mono-mp3",
	OutputFormatWebM:  "webm-24khz-16bit-mono-opus",
}

// formatMetadata 格式对应的 MIME type 和文件扩展名
var formatMetadata = map[string]struct {
	ContentType  string
	Extension    string
}{
	OutputFormatMP3:   {"audio/mpeg", ".mp3"},
	OutputFormatMP3HQ: {"audio/mpeg", ".mp3"},
	OutputFormatWebM:  {"audio/webm", ".webm"},
}

// validateOutputFormat 校验格式字符串是否在白名单中
func validateOutputFormat(format string) (string, error) {
	edgeFormat, ok := edgeOutputFormats[format]
	if !ok {
		return "", fmt.Errorf("不支持的输出格式 %q，可选值: %s", format, supportedFormatList())
	}
	return edgeFormat, nil
}

// supportedFormatList 返回合法格式列表（用于错误提示）
func supportedFormatList() string {
	list := make([]string, 0, len(edgeOutputFormats))
	for k := range edgeOutputFormats {
		list = append(list, k)
	}
	return fmt.Sprintf("%v", list)
}

// SetOutputFormat 设置输出格式
func SetOutputFormat(format string) CommunicateOption {
	return func(c *Communicate) error {
		edgeFormat, err := validateOutputFormat(format)
		if err != nil {
			return err
		}
		c.outputFormat = format
		_ = edgeFormat // 存储原始格式键，需要时通过 map 转换
		return nil
	}
}

// GetContentType 返回当前输出格式对应的 MIME 类型
func (c *Communicate) GetContentType() string {
	if meta, ok := formatMetadata[c.outputFormat]; ok {
		return meta.ContentType
	}
	return "audio/mpeg"
}

// GetFileExtension 返回当前输出格式对应的文件扩展名
func (c *Communicate) GetFileExtension() string {
	if meta, ok := formatMetadata[c.outputFormat]; ok {
		return meta.Extension
	}
	return ".mp3"
}
