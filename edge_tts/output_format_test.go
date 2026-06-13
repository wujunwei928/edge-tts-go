package edge_tts

import (
	"strings"
	"testing"
)

func TestSetOutputFormat_RejectsInvalidFormat(t *testing.T) {
	_, err := NewCommunicate("测试", SetOutputFormat("invalid"))
	if err == nil {
		t.Fatal("期望返回错误，但得到了 nil")
	}
}

func TestSetOutputFormat_AcceptsValidFormats(t *testing.T) {
	formats := []string{OutputFormatMP3, OutputFormatMP3HQ, OutputFormatWebM}
	for _, format := range formats {
		_, err := NewCommunicate("测试", SetOutputFormat(format))
		if err != nil {
			t.Fatalf("格式 %q 应该合法，但返回错误: %v", format, err)
		}
	}
}

func TestDefaultOutputFormat_IsMP3(t *testing.T) {
	c, err := NewCommunicate("测试")
	if err != nil {
		t.Fatalf("创建 Communicate 失败: %v", err)
	}
	cmd := c.getCommandRequestContent()
	if !strings.Contains(cmd, `"outputFormat":"audio-24khz-48kbitrate-mono-mp3"`) {
		t.Fatalf("默认输出应包含 mp3 格式，实际内容:\n%s", cmd)
	}
}

func TestSetOutputFormat_Passthrough(t *testing.T) {
	tests := []struct {
		format      string
		expectEdge  string
	}{
		{OutputFormatMP3, "audio-24khz-48kbitrate-mono-mp3"},
		{OutputFormatMP3HQ, "audio-24khz-96kbitrate-mono-mp3"},
		{OutputFormatWebM, "webm-24khz-16bit-mono-opus"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			c, err := NewCommunicate("测试", SetOutputFormat(tt.format))
			if err != nil {
				t.Fatalf("创建 Communicate 失败: %v", err)
			}
			cmd := c.getCommandRequestContent()
			expected := `"outputFormat":"` + tt.expectEdge + `"`
			if !strings.Contains(cmd, expected) {
				t.Fatalf("格式 %q 期望包含 %q，实际内容:\n%s", tt.format, expected, cmd)
			}
		})
	}
}

func TestGetContentType(t *testing.T) {
	tests := []struct {
		format   string
		expected string
	}{
		{OutputFormatMP3, "audio/mpeg"},
		{OutputFormatMP3HQ, "audio/mpeg"},
		{OutputFormatWebM, "audio/webm"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			c, _ := NewCommunicate("测试", SetOutputFormat(tt.format))
			if got := c.GetContentType(); got != tt.expected {
				t.Fatalf("GetContentType() = %q, 期望 %q", got, tt.expected)
			}
		})
	}
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		format   string
		expected string
	}{
		{OutputFormatMP3, ".mp3"},
		{OutputFormatMP3HQ, ".mp3"},
		{OutputFormatWebM, ".webm"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			c, _ := NewCommunicate("测试", SetOutputFormat(tt.format))
			if got := c.GetFileExtension(); got != tt.expected {
				t.Fatalf("GetFileExtension() = %q, 期望 %q", got, tt.expected)
			}
		})
	}
}

func TestDefaultMetadata(t *testing.T) {
	c, _ := NewCommunicate("测试")
	if got := c.GetContentType(); got != "audio/mpeg" {
		t.Fatalf("默认 GetContentType() = %q, 期望 %q", got, "audio/mpeg")
	}
	if got := c.GetFileExtension(); got != ".mp3" {
		t.Fatalf("默认 GetFileExtension() = %q, 期望 %q", got, ".mp3")
	}
}
