# edge-tts-go Web 子命令设计

## 概述

为 edge-tts-go 新增 `web` 子命令，提供本地 Web UI 进行文本转语音操作。定位为本地开发/调试工具，单二进制零额外依赖。

## 使用方式

```bash
edge-tts-go web [--port 8080] [--proxy ""]
```

- `--port`：监听端口，默认 `8080`
- `--proxy`：代理地址，透传给 edge_tts 包

## 架构方案

选择方案 A：单二进制 + `embed.FS` 内嵌 HTML/JS/CSS。

- 零额外依赖，保持项目现有依赖不变
- 单二进制分发，`edge-tts-go web` 即可启动
- 前端用 localStorage 存历史，无后端状态
- 与现有 CLI 模式完全解耦

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/` | 返回内嵌的 HTML 页面 |
| GET | `/api/voices` | 返回可用语音列表（JSON） |
| POST | `/api/tts` | 接收参数，返回 MP3 音频流 |

### POST /api/tts

请求体：

```json
{
  "text": "你好世界",
  "voice": "zh-CN-XiaoxiaoNeural",
  "rate": "+0%",
  "volume": "+0%",
  "pitch": "+0Hz"
}
```

响应：`audio/mpeg` 二进制流。

### GET /api/voices

响应：`edge_tts.Voice` 数组的 JSON。

## 文件结构

### 新增文件

```
edge-tts-go/
├── internal/
│   ├── cmd/
│   │   ├── root.go          # 修改：注册 web 子命令
│   │   └── web.go           # 新增：web 子命令 + HTTP handler
│   └── static/
│       └── index.html        # 新增：内嵌的 Web UI（HTML + CSS + JS 一体）
```

### 文件职责

- **`internal/cmd/web.go`**（约 120 行）：定义 Cobra 子命令、HTTP handler、使用 `//go:embed` 内嵌 HTML、启动 HTTP server
- **`internal/static/index.html`**（约 300-400 行）：单文件包含 HTML 结构 + `<style>` + `<script>`，纯原生 JS，使用 `localStorage` 存储历史记录

### edge_tts 包改动

无需修改。现有 `Communicate.Stream()` 返回 `[]byte`，完全满足 API 需求。

## Web UI 设计

### 页面布局

- 顶部标题栏：`edge-tts-go Web TTS`
- 文本输入区：`<textarea>` 多行输入
- 参数配置区：
  - 语音选择器（`<select>`），页面加载时通过 `/api/voices` 填充，按 Locale 分组
  - 语速滑块：`-100%` ~ `+300%`，默认 `+0%`
  - 音量滑块：`-100%` ~ `+100%`，默认 `+0%`
  - 音调滑块：`-50Hz` ~ `+50Hz`，默认 `+0Hz`
- 生成按钮：点击后 loading，完成后自动播放
- 历史记录区：显示之前的配音记录，支持重新播放和删除

### 历史记录

存储在浏览器 `localStorage`，每条记录包含：

```json
{
  "text": "你好世界",
  "voice": "zh-CN-XiaoxiaoNeural",
  "rate": "+0%",
  "volume": "+0%",
  "pitch": "+0Hz",
  "timestamp": 1710758400000,
  "audioData": "base64编码的mp3音频"
}
```

### 样式

纯 CSS，不引入框架。简洁风格，响应式布局适配桌面和移动端。

## 功能范围

包含：
- 基本配音功能（文本输入、语音选择、参数调节、生成播放下载）
- 历史记录（localStorage 存储、播放、删除）

不包含：
- 实时预览
- 文件上传
- 用户认证
- 多用户并发
