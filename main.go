package main

import (
	"fmt"
	"log"
	"os"

	"github.com/seven4x/videosrt/app"
	"github.com/spf13/cobra"
)

var (
	fileName string
	instance = app.NewByConfig(CONFIG)
)

//定义配置文件
const CONFIG = "config-prod.ini"

var rootCmd = &cobra.Command{
	Use:   "srt",
	Short: "srt",
	Long:  `srt`,
	Run: func(cmd *cobra.Command, args []string) {
		if fileName == "" {
			println("缺少参数")
			return
		}
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

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "clear",
	Long:  `clear`,
	Run: func(cmd *cobra.Command, args []string) {
		instance.ClearOssFile(app.GetMp3FileList(app.Mp3UploadRecord))
		log.Print("清除成功")
		err := os.Remove(app.Mp3UploadRecord)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&fileName, "fileName", "f", "demo.mp4", "fileName")
	mp3Cmd.PersistentFlags().StringVarP(&fileName, "fileName", "f", "demo.mp4", "fileName")
	_ = rootCmd.MarkFlagRequired("fileName")
	_ = mp3Cmd.MarkFlagRequired("fileName")
	rootCmd.AddCommand(mp3Cmd)
	rootCmd.AddCommand(clearCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func main() {
	Execute()
}
