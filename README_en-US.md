# edge-tts-go

* [中文](https://github.com/wujunwei928/edge-tts-go/blob/main/README.md)
* [English](https://github.com/wujunwei928/edge-tts-go/blob/main/README_en-US.md)

`edge-tts-go` is a golang module that allows you to use Microsoft Edge's online text-to-speech service from within your golang code or using the provided `edge-tts-go` command.

## Installation

### go install
To install it, run the following command:

    $ go install github.com/wujunwei928/edge-tts-go

### download release
    https://github.com/wujunwei928/edge-tts-go/releases

## Usage

### Basic usage

If you want to use the `edge-tts-go` command, you can simply run it with the following command:

    $ edge-tts-go --text "Hello, world" --write-media hello.mp3

### Changing the voice

If you want to change the language of the speech or more generally, the voice. 

You must first check the available voices with the `--list-voices` option:

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

    $ edge-tts-go --voice zh-CN-XiaoxiaoNeural --text "纵使浮云蔽天日，我亦拔剑破长空" --write-media hello_in_chinese.mp3

    if you already install ffplay in your computer, you can play the audio file with the following command:
    $ edge-tts-go --voice zh-CN-XiaoxiaoNeural --text "纵使浮云蔽天日，我亦拔剑破长空" | ffplay -i -
### Changing rate, volume and pitch

It is possible to make minor changes to the generated speech.

    $ edge-tts-go --rate=-50% --text "Hello, world" --write-media hello_with_rate_halved.mp3
    $ edge-tts-go --volume=-50% --text "Hello, world" --write-media hello_with_volume_halved.mp3
    $ edge-tts-go --pitch=-50Hz --text "Hello, world" --write-media hello_with_pitch_halved.mp3

## go module

It is possible to use the `edge-tts-go` module directly from go. look the following file:

* https://github.com/wujunwei928/edge-tts-go/blob/main/internal/cmd/root.go

## thanks

* https://github.com/rany2/edge-tts
