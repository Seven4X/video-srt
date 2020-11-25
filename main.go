package main

import (
	"fmt"
	"github.com/seven4x/videosrt/app"
	"github.com/spf13/cobra"
	"os"
)

var (
	fileName string
	instance = app.NewByConfig(CONFIG)
)

//定义配置文件
const CONFIG = "config.ini"

var rootCmd = &cobra.Command{
	Use:   "srt",
	Short: "srt",
	Long:  `srt`,
	Run: func(cmd *cobra.Command, args []string) {
		instance.Run(fileName)
	},
}
var mp3Cmd = &cobra.Command{
	Use:   "audio",
	Short: "audio",
	Long:  `audio`,
	Run: func(cmd *cobra.Command, args []string) {
		instance.RunWithoutExtract(fileName)
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&fileName, "fileName", "f", "demo.mp4", "fileName")
	rootCmd.MarkFlagRequired("fileName")
	rootCmd.AddCommand(mp3Cmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func main() {
	Execute()
}
