package ffmpeg

import (
	"github.com/floostack/transcoder/ffmpeg"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

var (
	ffmpegConf = &ffmpeg.Config{
		FfmpegBinPath:   "/usr/local/bin/ffmpeg",
		FfprobeBinPath:  "/usr/local/bin/ffprobe",
		ProgressEnabled: true,
	}
)

func init() {
	println("当前系统:\t" + runtime.GOOS)
	if runtime.GOOS == "windows" {
		ffmpegConf.FfmpegBinPath = "./bin/win/ffmpeg"
		ffmpegConf.FfprobeBinPath = "./bin/win/ffprobe"
	}
}
func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//提取视频音频
func ExtractAudio(video string, audioType string) (string, error) {
	overwrite := true
	format := "mp3"
	//aliyun目前支持最高采样率
	audioRate := 16000
	filename := strings.TrimSuffix(video, path.Ext(video))
	output := filename + audioType

	if exist(output) {
		log.Print("已有mp3文件，无需提取")
		return output, nil
	}
	opts := ffmpeg.Options{
		OutputFormat: &format,
		Overwrite:    &overwrite,
		AudioRate:    &audioRate,
	}
	progress, err := ffmpeg.
		New(ffmpegConf).
		Input(video).
		Output(output).
		WithOptions(opts).
		Start(opts)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	for msg := range progress {
		log.Printf("%+v", msg)
	}

	return output, nil
}
