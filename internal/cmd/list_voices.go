package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wujunwei928/edge-tts-go/edge_tts"
	"strings"
)

func NewCmdListVoices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-voices",
		Short: "lists available voices and exits",
		Long:  `lists available voices and exits`,
		RunE: func(cmd *cobra.Command, args []string) error {
			voiceList, err := edge_tts.ListVoices(proxyURL)
			if err != nil {
				return err
			}

			for _, voice := range voiceList {
				fmt.Println("Name:", voice.Name)
				fmt.Println("ShortName:", voice.ShortName)
				fmt.Println("Gender:", voice.Gender)
				fmt.Println("Locale:", voice.Locale)
				fmt.Println("ContentCategories:", strings.Join(voice.VoiceTag.ContentCategories, ","))
				fmt.Println("VoicePersonalities:", strings.Join(voice.VoiceTag.VoicePersonalities, ","))
				fmt.Println()
			}

			return nil
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(NewCmdListVoices())
}
