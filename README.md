# edge-tts-go

* [中文](https://github.com/wujunwei928/edge-tts-go/blob/main/README.md)
* [English](https://github.com/wujunwei928/edge-tts-go/blob/main/README_en-US.md)

`edge-tts-go` 是一个 golang 模块，允许您从 golang 代码中或使用提供的 `edge-tts-go` 命令使用 Microsoft Edge 的在线文本到语音服务。

## 安装

### go install
    $ go install github.com/wujunwei928/edge-tts-go

### 下载预编译版本
    https://github.com/wujunwei928/edge-tts-go/releases

## 用法

### 基本用法

如果您想使用 `edge-tts-go` 命令，只需使用以下命令运行它：

    $ edge-tts-go --text "Hello, world" --write-media hello.mp3

### 改变声音

如果您想更改转换文本时使用的声音。

您需要使用 `--list-voices` 选项检查可用的语音：

    $ edge-tts-go --list-voices
    Name: Microsoft Server Speech Text to Speech Voice (zh-CN, XiaoxiaoNeural)
    ShortName: zh-CN-XiaoxiaoNeural
    Gender: Female
    Locale: zh-CN
    ContentCategories: News,Novel
    VoicePersonalities: Warm
    
    Name: Microsoft Server Speech Text to Speech Voice (zh-CN, XiaoyiNeural)
    ShortName: zh-CN-XiaoyiNeural
    Gender: Female
    Locale: zh-CN
    ContentCategories: Cartoon,Novel
    VoicePersonalities: Lively

    ...

使用 `--voice` 选项指定声音进行转换

    $ edge-tts-go --voice zh-CN-XiaoxiaoNeural --text "纵使浮云蔽天日，我亦拔剑破长空" --write-media hello_in_chinese.mp3

    如果你的电脑安装过ffplay，你可以使用以下命令直接播放音频文件:
    $ edge-tts-go --voice zh-CN-XiaoxiaoNeural --text "纵使浮云蔽天日，我亦拔剑破长空" | ffplay -i -

### 改变速率、音量和音调

    $ edge-tts-go --rate=-50% --text "Hello, world" --write-media hello_with_rate_halved.mp3
    $ edge-tts-go --volume=-50% --text "Hello, world" --write-media hello_with_volume_halved.mp3
    $ edge-tts-go --pitch=-50Hz --text "Hello, world" --write-media hello_with_pitch_halved.mp3

### 改变输出格式

支持三种输出格式，默认为 MP3（48kbps）：

| 格式 | 说明 |
|---|---|
| `mp3` | MP3 48kbps（默认） |
| `mp3-hq` | MP3 96kbps，更高质量 |
| `webm` | WebM Opus |

通过 `--format` 指定：

    $ edge-tts-go --text "Hello, world" --format mp3-hq --write-media hello_hq.mp3
    $ edge-tts-go --text "Hello, world" --format webm --write-media hello.webm

也可以从 `--write-media` 文件扩展名自动推断（仅 `.webm`）：

    $ edge-tts-go --text "Hello, world" --write-media hello.webm

## go 模块

可以直接在go代码中使用 `edge-tts-go` 模块：

```go
package main

import (
    "os"
    "github.com/wujunwei928/edge-tts-go/edge_tts"
)

func main() {
    c, err := edge_tts.NewCommunicate("你好世界",
        edge_tts.SetVoice("zh-CN-XiaoxiaoNeural"),
        edge_tts.SetOutputFormat(edge_tts.OutputFormatMP3HQ), // 可选：mp3 / mp3-hq / webm
    )
    if err != nil {
        panic(err)
    }

    audioData, err := c.Stream()
    if err != nil {
        panic(err)
    }

    os.WriteFile("output"+c.GetFileExtension(), audioData, 0644)
}
```

## 致谢

* https://github.com/rany2/edge-tts
