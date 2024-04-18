package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wujunwei928/edge-tts-go/edge_tts"
	"io"
	"os"
)

func NewCmdTextToSpeech() *cobra.Command {
	var (
		text       string
		file       string
		voice      string
		rate       string
		volume     string
		pitch      string
		wordsInCue float64
		writeMedia string
	)

	cmd := &cobra.Command{
		Use:   "text-to-speech",
		Short: "text to speech",
		Long:  `text to speech`,
		RunE: func(cmd *cobra.Command, args []string) error {
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
			fmt.Println("text", inputText)
			fmt.Println("voice", voice)

			conn, err := edge_tts.NewCommunicate(
				inputText,
				edge_tts.SetVoice(voice),
			)
			fmt.Println(conn)
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

			// 将mp3文件的二进制数据写入标准输出
			_, err = io.Copy(os.Stdout, bytes.NewReader(audioData))
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&text, "text", "t", "", "what TTS will say")
	cmd.Flags().StringVarP(&file, "file", "f", "", "same as --text but read from file, when val is /dev/stdin, read from stdin")
	cmd.Flags().StringVarP(&voice, "voice", "v", "en-US-AriaNeural", "voice for TTS")
	cmd.Flags().StringVar(&rate, "rate", "+0%", "set TTS rate")
	cmd.Flags().StringVar(&volume, "volume", "+0%", "set TTS volume")
	cmd.Flags().StringVar(&pitch, "pitch", "+0Hz", "set TTS pitch")
	cmd.Flags().Float64Var(&wordsInCue, "words-in-cue", 10, "number of words in a subtitle cue")
	cmd.Flags().StringVar(&writeMedia, "write-media", "", "send media output to file instead of stdout")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewCmdTextToSpeech())
}
