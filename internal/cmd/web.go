package cmd

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/wujunwei928/edge-tts-go/edge_tts"
)

//go:embed static/index.html
var staticFS embed.FS

var (
	webPort  int
	webProxy string
)

// ttsRequest 定义 TTS 请求参数
type ttsRequest struct {
	Text   string `json:"text"`
	Voice  string `json:"voice"`
	Rate   string `json:"rate"`
	Volume string `json:"volume"`
	Pitch  string `json:"pitch"`
}

// handleIndex 返回内嵌的 HTML 页面
func handleIndex(w http.ResponseWriter, r *http.Request) {
	data, err := staticFS.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "页面加载失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

// handleVoices 返回可用语音列表
func handleVoices(w http.ResponseWriter, r *http.Request) {
	voices, err := edge_tts.ListVoices(webProxy)
	if err != nil {
		http.Error(w, fmt.Sprintf("获取语音列表失败: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(voices)
}

// handleTTS 接收文本参数，返回 MP3 音频流
func handleTTS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持 POST", http.StatusMethodNotAllowed)
		return
	}

	var req ttsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求参数解析失败", http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		http.Error(w, "文本不能为空", http.StatusBadRequest)
		return
	}

	options := []edge_tts.CommunicateOption{
		edge_tts.SetVoice(req.Voice),
		edge_tts.SetRate(req.Rate),
		edge_tts.SetVolume(req.Volume),
		edge_tts.SetPitch(req.Pitch),
		edge_tts.SetReceiveTimeout(30),
	}
	if webProxy != "" {
		options = append(options, edge_tts.SetProxy(webProxy))
	}

	conn, err := edge_tts.NewCommunicate(req.Text, options...)
	if err != nil {
		http.Error(w, fmt.Sprintf("创建连接失败: %v", err), http.StatusInternalServerError)
		return
	}

	audioData, err := conn.Stream()
	if err != nil {
		http.Error(w, fmt.Sprintf("语音合成失败: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Disposition", "inline; filename=\"tts.mp3\"")
	w.Write(audioData)
}

// webCmd 定义 web 子命令
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "启动 Web UI 进行文本配音",
	Long:  `启动本地 Web 服务，通过浏览器进行文本转语音操作。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mux := http.NewServeMux()
		mux.HandleFunc("/", handleIndex)
		mux.HandleFunc("/api/voices", handleVoices)
		mux.HandleFunc("/api/tts", handleTTS)

		addr := fmt.Sprintf(":%d", webPort)
		log.Printf("Web UI 已启动: http://localhost:%d", webPort)
		return http.ListenAndServe(addr, mux)
	},
}

func init() {
	webCmd.Flags().IntVar(&webPort, "port", 8080, "Web 服务监听端口")
	webCmd.Flags().StringVar(&webProxy, "proxy", "", "使用代理访问 TTS 服务")
}
