# ADR 0001: 输出格式支持与设计

## 状态

已接受

## 日期

2026-06-13

## 背景

用户希望 edge-tts-go 支持多种音频输出格式，特别是 WAV（`riff-24khz-16bit-mono-pcm`）。当前代码在 `communicate.go` 的 `getCommandRequestContent()` 中硬编码了 `audio-24khz-48kbitrate-mono-mp3`，没有任何配置项。

## 调研结论

通过编写验证脚本，对 Edge TTS 免费 WebSocket 端点（`speech.platform.bing.com`）测试了 Azure Speech SDK 定义的全部 38 种输出格式，**仅 3 种可用**：

| 格式 | 状态 | 数据量（"你好"） |
|---|---|---|
| `audio-24khz-48kbitrate-mono-mp3` | ✅ | 7,200 bytes |
| `audio-24khz-96kbitrate-mono-mp3` | ✅ | 14,400 bytes |
| `webm-24khz-16bit-mono-opus` | ✅ | 7,710 bytes |
| 其他 35 种（含所有 RIFF/PCM/OGG/AMR） | ❌ | 0 bytes |

这与 Python edge-tts 维护者在 [issue #392](https://github.com/rany2/edge-tts/issues/392) 中的结论一致："only webm and mp3 formats are supported by the service"。

## 决策

### 1. 不引入客户端转码

**选择**：保持薄封装原则，不引入 MP3→WAV 转码依赖。

**理由**：
- `edge_tts/` 包的定位是 Edge TTS WebSocket 协议的薄封装，引入音频编解码会让它变成"音频处理库"
- 会增加外部依赖（MP3 解码库、WAV 编码库）和二进制体积
- 用户需要 WAV 时可在下游自行转码，一行代码即可完成

**替代方案**：在库内或独立子包 `edge_tts/convert` 中提供转码。因偏离简单原则而否决。

### 2. 单层 API + 预定义常量

**选择**：`SetOutputFormat(format string)` 接受预定义常量值，对用户屏蔽 Edge 内部格式细节。

**理由**：
- 用户说"我要 MP3"就够了，不需要知道 Edge 用 48k 还是 96k
- 常量名语义化：`OutputFormatMP3`、`OutputFormatMP3HQ`、`OutputFormatWebM`
- 白名单严格校验，传入非法值立即返回 error

**替代方案**：暴露 Edge 原始格式字符串（如 `audio-24khz-48kbitrate-mono-mp3`）。因字符串过长、用户认知负担大而否决。

### 3. 白名单仅 3 个值

**选择**：严格白名单，只有实测通过的 3 种格式。

**理由**：
- 格式字符串拼错会导致 API 返回难以理解的错误，白名单能在最上层给出清晰提示
- Edge TTS 格式集稳定且有限，维护成本极低

### 4. 消费端自适应

**选择**：`Communicate` 暴露 `GetContentType()` 和 `GetFileExtension()` 方法，CLI/Web 消费端根据实际格式动态设置响应头。

**理由**：格式→MIME/扩展名的映射是 `edge_tts` 包的领域知识，放在这里最内聚，且不改变 `Stream()` 的现有签名。

## 影响

- 新增文件 `edge_tts/output_format.go`（常量、白名单、映射、校验）
- 修改 `edge_tts/communicate.go`（结构体字段、Option 函数、元数据方法）
- 修改 `internal/cmd/root.go`（`--format` flag + 扩展名推断）
- 修改 `internal/cmd/web.go`（format 参数 + 动态 Content-Type）
- 修改 `internal/cmd/static/index.html`（Web UI 格式选择器）
