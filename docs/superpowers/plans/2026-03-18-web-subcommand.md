# Web 子命令实施计划

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为 edge-tts-go 新增 `web` 子命令，提供本地 Web UI 进行文本配音。

**Architecture:** 单二进制 + `embed.FS` 内嵌 HTML。Go 标准库 `net/http` 提供 3 个 API 端点（静态页面、语音列表、TTS 转换）。前端单文件 HTML/CSS/JS，使用 `localStorage` 存储历史记录。不修改 `edge_tts` 核心包。

**Tech Stack:** Go 标准库 `net/http`、`embed`、`encoding/json`；纯原生 HTML/CSS/JS

---

## 文件变更总览

| 文件 | 操作 | 职责 |
|------|------|------|
| `internal/cmd/web.go` | 新建 | web 子命令定义、HTTP handler、embed 指令、server 启动 |
| `internal/static/index.html` | 新建 | Web UI 页面（HTML + CSS + JS 一体） |
| `internal/cmd/root.go` | 修改 | 在 `init()` 中注册 web 子命令 |

---

### Task 1: 创建 web.go — HTTP handler 与子命令骨架

**Files:**
- Create: `internal/cmd/web.go`

- [ ] **Step 1: 创建 `internal/cmd/web.go`，实现 embed 指令和 HTTP handler**

```go
package cmd

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
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
		log.Printf("🎉 Web UI 已启动: http://localhost:%d", webPort)
		return http.ListenAndServe(addr, mux)
	},
}

func init() {
	webCmd.Flags().IntVar(&webPort, "port", 8080, "Web 服务监听端口")
	webCmd.Flags().StringVar(&webProxy, "proxy", "", "使用代理访问 TTS 服务")
}
```

- [ ] **Step 2: 验证编译通过**

Run: `cd /code/edge-tts-go && go build ./...`
Expected: 编译成功（此时 `internal/static/index.html` 不存在会报 embed 错误，需先创建空文件）

- [ ] **Step 3: 创建空的 `internal/static/index.html` 占位**

```bash
mkdir -p /code/edge-tts-go/internal/static
touch /code/edge-tts-go/internal/static/index.html
```

- [ ] **Step 4: 再次验证编译通过**

Run: `cd /code/edge-tts-go && go build ./...`
Expected: 编译成功

- [ ] **Step 5: 提交**

```bash
git add internal/cmd/web.go internal/static/index.html
git commit -m "feat: 添加 web 子命令骨架和 HTTP handler"
```

---

### Task 2: 注册 web 子命令到 rootCmd

**Files:**
- Modify: `internal/cmd/root.go:109-115`

- [ ] **Step 1: 在 `root.go` 的 `Execute()` 函数中注册 web 子命令**

在 `func init()` 末尾添加：

```go
rootCmd.AddCommand(webCmd)
```

修改后 `root.go` 的 `init()` 函数应为：

```go
func init() {
	// bind flags
	rootCmd.Flags().StringVarP(&text, "text", "t", "", "what TTS will say")
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "same as --text but read from file")
	rootCmd.Flags().StringVarP(&voice, "voice", "v", "en-US-AriaNeural", "voice for TTS")
	rootCmd.Flags().StringVar(&rate, "rate", "+0%", "set TTS rate")
	rootCmd.Flags().StringVar(&volume, "volume", "+0%", "set TTS volume")
	rootCmd.Flags().StringVar(&pitch, "pitch", "+0Hz", "set TTS pitch")
	rootCmd.Flags().Float64Var(&wordsInCue, "words-in-cue", 10, "number of words in a subtitle cue")
	rootCmd.Flags().StringVar(&writeMedia, "write-media", "", "send media output to file instead of stdout")
	rootCmd.Flags().StringVar(&proxyURL, "proxy", "", "use a proxy for TTS and voice list")
	rootCmd.Flags().BoolVar(&listVoices, "list-voices", false, "lists available voices and exits")

	rootCmd.AddCommand(webCmd)
}
```

- [ ] **Step 2: 验证编译并检查子命令注册**

Run: `cd /code/edge-tts-go && go build -o edge-tts-go && ./edge-tts-go --help`
Expected: 输出中包含 `web` 子命令

- [ ] **Step 3: 提交**

```bash
git add internal/cmd/root.go
git commit -m "feat: 注册 web 子命令到 rootCmd"
```

---

### Task 3: 实现 Web UI — index.html 完整页面

**Files:**
- Modify: `internal/static/index.html`

这是最大的单个任务。HTML 文件包含以下部分：

- [ ] **Step 1: 编写完整的 `internal/static/index.html`**

文件结构分为三块：`<style>` CSS、HTML 结构、`<script>` JS。

**CSS 部分要点：**
- CSS 变量定义主题色
- `.container` 居中布局，最大宽度 720px
- `.form-group` 控件间距
- `textarea` 宽高自适应
- `input[type=range]` 自定义滑块样式
- `.btn-primary` 生成按钮，蓝色调
- `.btn-primary:disabled` loading 态
- `.history-item` 历史记录卡片布局
- `.spinner` 加载动画
- 响应式：`@media (max-width: 600px)` 适配移动端

**HTML 结构要点：**
- `<header>` 标题栏
- `<main>` 包含：
  - `<textarea id="textInput">` 文本输入
  - 语音选择 `<select id="voiceSelect">` + loading 状态
  - 语速滑块 `<input type="range" id="rateSlider">` + 显示值 `<span>`
  - 音量滑块 `<input type="range" id="volumeSlider">` + 显示值
  - 音调滑块 `<input type="range" id="pitchSlider">` + 显示值
  - `<button id="generateBtn">` 生成按钮
  - `<audio id="audioPlayer" controls>` 隐藏的音频播放器
  - `<div id="errorMsg">` 错误提示区
  - `<section id="historySection">` 历史记录区
- `<footer>` 版权信息

**JS 部分要点：**

```javascript
// 常量
const HISTORY_KEY = 'edge-tts-history';
const MAX_HISTORY = 50;

// DOM 元素引用
const elements = {
  textInput: document.getElementById('textInput'),
  voiceSelect: document.getElementById('voiceSelect'),
  rateSlider: document.getElementById('rateSlider'),
  volumeSlider: document.getElementById('volumeSlider'),
  pitchSlider: document.getElementById('pitchSlider'),
  rateValue: document.getElementById('rateValue'),
  volumeValue: document.getElementById('volumeValue'),
  pitchValue: document.getElementById('pitchValue'),
  generateBtn: document.getElementById('generateBtn'),
  audioPlayer: document.getElementById('audioPlayer'),
  errorMsg: document.getElementById('errorMsg'),
  historyList: document.getElementById('historyList'),
};

// 页面加载时：获取语音列表 + 加载历史记录
async function init() {
  await loadVoices();
  loadHistory();
  bindEvents();
}

// 加载语音列表，按 Locale 分组填充 <select>
async function loadVoices() {
  const resp = await fetch('/api/voices');
  const voices = await resp.json();
  // 按 Locale 分组
  const grouped = {};
  voices.forEach(v => {
    if (!grouped[v.Locale]) grouped[v.Locale] = [];
    grouped[v.Locale].push(v);
  });
  // 填充 optgroup
  elements.voiceSelect.innerHTML = '';
  Object.keys(grouped).sort().forEach(locale => {
    const group = document.createElement('optgroup');
    group.label = locale;
    grouped[locale].forEach(v => {
      const opt = document.createElement('option');
      opt.value = v.ShortName;
      opt.textContent = `${v.ShortName} (${v.Gender})`;
      group.appendChild(opt);
    });
    elements.voiceSelect.appendChild(group);
  });
}

// 滑块值格式化
function formatRate(val) { return `${val >= 0 ? '+' : ''}${val}%`; }
function formatVolume(val) { return `${val >= 0 ? '+' : ''}${val}%`; }
function formatPitch(val) { return `${val >= 0 ? '+' : ''}${val}Hz`; }

// 生成 TTS
async function generate() {
  const text = elements.textInput.value.trim();
  if (!text) return showError('请输入文本');

  setGenerating(true);
  try {
    const resp = await fetch('/api/tts', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        text,
        voice: elements.voiceSelect.value,
        rate: formatRate(elements.rateSlider.value),
        volume: formatVolume(elements.volumeSlider.value),
        pitch: formatPitch(elements.pitchSlider.value),
      }),
    });

    if (!resp.ok) {
      const errText = await resp.text();
      throw new Error(errText);
    }

    const blob = await resp.blob();
    const url = URL.createObjectURL(blob);
    elements.audioPlayer.src = url;
    elements.audioPlayer.play();
    saveHistory(text, blob);
  } catch (e) {
    showError(e.message);
  } finally {
    setGenerating(false);
  }
}

// localStorage 历史管理
function saveHistory(text, blob) {
  const reader = new FileReader();
  reader.onload = () => {
    const history = JSON.parse(localStorage.getItem(HISTORY_KEY) || '[]');
    const record = {
      id: Date.now(),
      text,
      voice: elements.voiceSelect.value,
      rate: formatRate(elements.rateSlider.value),
      volume: formatVolume(elements.volumeSlider.value),
      pitch: formatPitch(elements.pitchSlider.value),
      timestamp: Date.now(),
      audioData: reader.result, // base64
    };
    history.unshift(record);
    if (history.length > MAX_HISTORY) history.pop();
    localStorage.setItem(HISTORY_KEY, JSON.stringify(history));
    renderHistory();
  };
  reader.readAsDataURL(blob);
}

function loadHistory() {
  renderHistory();
}

function deleteHistory(id) {
  let history = JSON.parse(localStorage.getItem(HISTORY_KEY) || '[]');
  history = history.filter(h => h.id !== id);
  localStorage.setItem(HISTORY_KEY, JSON.stringify(history));
  renderHistory();
}

function renderHistory() {
  const history = JSON.parse(localStorage.getItem(HISTORY_KEY) || '[]');
  // 渲染历史列表 DOM
  elements.historyList.innerHTML = history.length === 0
    ? '<p class="empty-hint">暂无历史记录</p>'
    : history.map(h => createHistoryItem(h)).join('');
  // 绑定播放/删除事件（事件委托）
}

function createHistoryItem(record) {
  return `
    <div class="history-item" data-id="${record.id}">
      <div class="history-info">
        <span class="history-text">${escapeHtml(record.text.slice(0, 50))}</span>
        <span class="history-meta">${record.voice} · ${record.rate} · ${record.volume} · ${record.pitch}</span>
      </div>
      <div class="history-actions">
        <button class="btn-play" onclick="playHistory(${record.id})" title="播放">▶</button>
        <button class="btn-delete" onclick="deleteHistory(${record.id})" title="删除">✕</button>
      </div>
    </div>`;
}

function playHistory(id) {
  const history = JSON.parse(localStorage.getItem(HISTORY_KEY) || '[]');
  const record = history.find(h => h.id === id);
  if (record && record.audioData) {
    elements.audioPlayer.src = record.audioData;
    elements.audioPlayer.play();
  }
}

// 工具函数
function escapeHtml(str) { ... }
function showError(msg) { ... }
function setGenerating(val) { ... }

// 事件绑定
function bindEvents() {
  elements.generateBtn.addEventListener('click', generate);
  // 滑块实时更新显示值
  elements.rateSlider.addEventListener('input', () => {
    elements.rateValue.textContent = formatRate(elements.rateSlider.value);
  });
  // ... volumeSlider, pitchSlider 同理
}

// 启动
init();
```

- [ ] **Step 2: 验证页面加载**

Run: `cd /code/edge-tts-go && go build -o edge-tts-go && ./edge-tts-go web &`
在浏览器访问 `http://localhost:8080`，确认页面正常渲染。

- [ ] **Step 3: 验证 TTS 功能**

在页面中输入文本、选择语音、点击生成，确认音频播放正常。

- [ ] **Step 4: 验证历史记录**

生成几次后刷新页面，确认历史记录持久化、播放和删除功能正常。

- [ ] **Step 5: 提交**

```bash
git add internal/static/index.html
git commit -m "feat: 实现 Web UI 页面（配音功能 + 历史记录）"
```

---

### Task 4: 集成测试与清理

**Files:**
- Modify: 无新文件变更

- [ ] **Step 1: 编译所有平台确认无问题**

Run: `cd /code/edge-tts-go && GOOS=linux GOARCH=amd64 go build -o edge-tts-go && GOOS=darwin GOARCH=amd64 go build -o edge-tts-go-darwin && GOOS=windows GOARCH=amd64 go build -o edge-tts-go.exe`
Expected: 全部编译成功

- [ ] **Step 2: 验证 `edge-tts-go web --help`**

Run: `./edge-tts-go web --help`
Expected: 显示 `--port` 和 `--proxy` flag 说明

- [ ] **Step 3: 验证原有 CLI 功能不受影响**

Run: `./edge-tts-go --list-voices`
Expected: 正常输出语音列表

- [ ] **Step 4: 最终提交**

如有任何修复，统一提交。

---

## 依赖检查

- **无新依赖**：仅使用 Go 标准库（`net/http`、`embed`、`encoding/json`、`fmt`、`log`、`io/fs`）
- **edge_tts 包无需修改**：直接使用现有 `ListVoices()` 和 `NewCommunicate().Stream()`
