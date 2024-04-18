package cmd

import (
	"bytes"
	"errors"
	"github.com/spf13/cobra"
	"github.com/wujunwei928/edge-tts-go/edge_tts"
	"io"
	"os"
)

var (
	listVoices bool
	text       string
	file       string
	voice      string
	rate       string
	volume     string
	pitch      string
	wordsInCue float64
	writeMedia string
	proxyURL   string // 是否使用代理
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "edge-tts-go",
	Short:   "调用Edge TTS服务，文本生成语音",
	Long:    `调用Edge TTS服务，文本生成语音`,
	Version: edge_tts.PackageVersion, // 指定版本号: 会有 -v 和 --version 选项, 用于打印版本号
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		// 列出可用的语音
		if listVoices {
			ListVoices()
			return nil
		}

		// 文本转语音

		if len(text) <= 0 && len(file) <= 0 {
			return errors.New("--text and --file can't be empty at the same time")
		}

		inputText := text
		if len(file) > 0 {
			var (
				fileContent []byte
				readFileErr error
			)

			switch file {
			case "/dev/stdin":
				fileContent, readFileErr = io.ReadAll(os.Stdin)
			default:
				fileContent, readFileErr = os.ReadFile(file)
			}

			if readFileErr != nil {
				return readFileErr
			}

			inputText = string(fileContent)
		}

		connOptions := []edge_tts.CommunicateOption{
			edge_tts.SetVoice(voice),
			edge_tts.SetRate(rate),
			edge_tts.SetVolume(volume),
			edge_tts.SetPitch(pitch),
			edge_tts.SetReceiveTimeout(20),
		}
		if len(proxyURL) > 0 {
			connOptions = append(connOptions, edge_tts.SetProxy(proxyURL))
		}

		conn, err := edge_tts.NewCommunicate(
			inputText,
			connOptions...,
		)
		if err != nil {
			return err
		}
		audioData, err := conn.Stream()
		if err != nil {
			return err
		}

		if len(writeMedia) > 0 {
			writeMediaErr := os.WriteFile(writeMedia, audioData, 0644)
			if writeMediaErr != nil {
				return writeMediaErr
			}
			return nil
		}

		// write mp3 file's binary data to stdout
		_, err = io.Copy(os.Stdout, bytes.NewReader(audioData))
		if err != nil {
			return err
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

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
}
