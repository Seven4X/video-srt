package main

import (
	"bufio"
	"fmt"
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
		instance.ClearOssFile(getMp3FileList(app.Mp3UploadRecord))
	},
}

func getMp3FileList(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		println(err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	mp3s := make([]string, 0)
	for scanner.Scan() {
		file := scanner.Text()
		mp3s = append(mp3s, file)
	}
	return mp3s
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
