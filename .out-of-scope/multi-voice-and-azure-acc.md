# 多发音人混合编排与 Azure ACC 免费接口

本项目的 TTS 能力基于 Microsoft Edge 浏览器「大声朗读」接口
（`wss://speech.platform.bing.com/.../readaloud/edge/v1`）。以下两类增强诉求
经实测评估后不予实现。

## 1. 单次请求内的多发音人混合 / 丰富 SSML

Edge ReadAloud 接口对 SSML 的接受范围非常窄。实测结论：

- ✅ 支持：恰好一个 `<voice>` + 其下任意多个 `<prosody>` 段（包裹纯文本）
- ❌ 拒绝：多个 `<voice>` 段（返回 `SSML is invalid`，0 字节音频）
- ❌ 拒绝：`<break>`、`<say-as>`、`<emphasis>`、`<sub>`、`<phoneme>`、
  `<p>`/`<s>`、`<mstts:express-as>`（情感风格）等 —— 连最基础的停顿和段落标签都被拒

因此「一份音频内不同句子用不同发音人」「情感风格、强调、停顿等丰富表达」
无法在现有接口上实现。要做到多发音人，只能由客户端侧「分段单 voice 合成 +
音频拼接」，属于调用方职责，不在核心库范围内。

## 2. 接入 Azure Speech Studio ACC 免费试用接口

`*.api.speech.microsoft.com/accfreetrial/.../vcg/speak`（即 speech.microsoft.com
网页端 Audio Content Creation 用的接口）经探测：

- 不需要显式鉴权 token 即可触达业务层，但 `ttsAudioFormat` 参数极挑剔 ——
  报告者使用的 `audio-24khz-160kbitrate-mono-mp3` 及另外试的 5 种常见 Azure
  格式全部返回 `400 InvalidPayload`，需从浏览器实际请求逆向出合法格式串；
- 限流严格（约 20/50 次每周期），且**失败的请求也消耗配额**；
- CORS 仅允许 `https://speech.microsoft.com` 来源。

该接口不稳定、配额小、参数需逆向、CORS 锁死官方站，不适合作为可长期依赖的
生产接口，故不接入。

## 现有可替代方案

- 「同文本多发音人」若指**批量生成 N 份独立音频文件**：现有 `--voice`
  （ShortName 格式）+ 循环调用 `Stream()` 已支持，无需改动。
- 「单发音人内分段差异化表达」：可在同一 `<voice>` 下用多个 `<prosody>` 段
  实现（接口已支持）。

## Prior requests

- #1 — 「同文本多发音人」+ 建议换用 Azure accfreetrial 接口
