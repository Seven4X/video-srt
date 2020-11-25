package ffmpeg

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/floostack/transcoder/ffmpeg"
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
	execPath := getExecPath()
	println("当前目录：\t" + execPath)
	if runtime.GOOS == "windows" {
		ffmpegConf.FfmpegBinPath = execPath + "/bin/win/ffmpeg"
		ffmpegConf.FfprobeBinPath = execPath + "/bin/win/ffprobe"
	}
}
func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func getExecPath() string {

	execPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}
	//    Is Symlink
	fi, err := os.Lstat(execPath)
	if err != nil {
		log.Fatal(err)
	}
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		execPath, err = os.Readlink(execPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	execDir := filepath.Dir(execPath)
	if execDir == "." {
		execDir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	return execDir
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
