# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

edge-tts-go 是一个 Go 语言实现的 Microsoft Edge TTS（文本转语音）客户端。支持 CLI 命令行和 Web UI 两种使用方式，也可作为 Go 模块被其他项目引用。

## 常用命令

```bash
# 构建
go build -o edge-tts-go

# 运行 CLI
./edge-tts-go --text "你好世界" --write-media hello.mp3
./edge-tts-go --list-voices

# 启动 Web UI
./edge-tts-go web --port 8080

# 静态检查与格式化
go vet ./...
go fmt ./...

# 运行指定测试文件
go test ./edge_tts/ -run TestFuncName -timeout 30s
```

注意：项目无正式测试套件，`test.go` 为 DRM 令牌的开发测试文件。

## 架构

### 核心分层

- `main.go` → 入口，调用 `internal/cmd.Execute()`
- `internal/cmd/` → Cobra CLI 层
  - `root.go` → 根命令（CLI 模式的 TTS 合成逻辑）
  - `web.go` → `web` 子命令（HTTP 服务器，嵌入式 Web UI）
  - `static/index.html` → 通过 `//go:embed` 嵌入的单页面 Web UI
- `edge_tts/` → 核心库（可被外部 Go 项目直接引用）
  - `communicate.go` → WebSocket 通信与 TTS 合成核心（`Communicate` 结构体，Option 模式配置）
  - `drm.go` → `Sec-MS-GEC` 令牌生成（SHA256 哈希 + 时钟偏差校正）
  - `constants.go` → API URL、Chromium 版本、HTTP Headers（模拟 Edge 浏览器）
  - `list_voices.go` → 语音列表获取（resty HTTP 客户端）
  - `eorrors.go` → 预定义错误

### 关键数据流

1. CLI/Web 入口 → 构建 `Communicate` 实例（通过 Option 函数：`SetVoice`, `SetRate`, `SetVolume`, `SetPitch`, `SetProxy`）
2. `Communicate.Stream()` → 建立 WebSocket 连接 → 发送配置 + SSML 请求 → 接收二进制音频数据（MP3）
3. 音频数据写入文件或通过 stdout 输出

### API 认证机制

- 使用 `Sec-MS-GEC` 令牌（SHA256 哈希，基于 Windows 文件时间 + TrustedClientToken）
- HTTP Headers 模拟 Edge 浏览器指纹（`CHROMIUM_FULL_VERSION` 需随 Edge 更新保持同步）
- Web API 端点：`/api/voices`（GET 语音列表）、`/api/tts`（POST 文本转语音）

## 关键约定

- Go 版本：1.21+
- CLI 参数 `--voice` 接受 ShortName 格式（如 `zh-CN-XiaoxiaoNeural`），内部会自动转换为完整名称
- Web UI 通过 `embed.FS` 嵌入，无需外部静态文件服务
- 无 Makefile，构建直接使用 `go build`
- CI 通过推送 `0.*` 格式 tag 触发跨平台构建（linux/windows/darwin amd64）
