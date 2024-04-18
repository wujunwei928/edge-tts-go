package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wujunwei928/edge-tts-go/edge_tts"
	"os"
)

var (
	proxyURL string // 是否使用代理
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "dev",
	Short:   "研发工具箱",
	Long:    `研发工具箱`,
	Version: edge_tts.PackageVersion, // 指定版本号: 会有 -v 和 --version 选项, 用于打印版本号
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	//
	rootCmd.PersistentFlags().StringVarP(&proxyURL, "proxy", "", "", "use a proxy for TTS and voice list.")
}
