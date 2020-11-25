package app

import (
	"bytes"
	"github.com/buger/jsonparser"
	"github.com/seven4x/videosrt/app/aliyun/cloud"
	"github.com/seven4x/videosrt/app/aliyun/oss"
	"github.com/seven4x/videosrt/app/ffmpeg"
	"github.com/seven4x/videosrt/lib/config/ini"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

//主应用
type VideoSrt struct {
	AliyunOss   oss.AliyunOss     //oss
	AliyunCloud cloud.AliyunCloud //语音识别引擎

	IntelligentBlock bool   //智能分段处理
	TempDir          string //临时文件目录
	AppDir           string //应用根目录
}

//读取配置
func NewByConfig(cfg string) *VideoSrt {
	if file, e := ini.LoadConfigFile(cfg, "."); e != nil {
		panic(e)
	} else {
		appconfig := &VideoSrt{}

		//AliyunOss
		appconfig.AliyunOss.Endpoint = file.GetMust("aliyunOss.endpoint", "")
		appconfig.AliyunOss.AccessKeyId = file.GetMust("aliyunOss.accessKeyId", "")
		appconfig.AliyunOss.AccessKeySecret = file.GetMust("aliyunOss.accessKeySecret", "")
		appconfig.AliyunOss.BucketName = file.GetMust("aliyunOss.bucketName", "")
		appconfig.AliyunOss.BucketDomain = file.GetMust("aliyunOss.bucketDomain", "")
		appconfig.AliyunOss.UploadPath = file.GetMust("aliyunOss.uploadPath", "")

		//AliyunCloud
		appconfig.AliyunCloud.AccessKeyId = file.GetMust("aliyunClound.accessKeyId", "")
		appconfig.AliyunCloud.AccessKeySecret = file.GetMust("aliyunClound.accessKeySecret", "")
		appconfig.AliyunCloud.AppKey = file.GetMust("aliyunClound.appKey", "")

		appconfig.IntelligentBlock = file.GetBoolMust("srt.intelligent_block", false)
		appconfig.TempDir = "temp/audio"

		return appconfig
	}
}

//应用运行
func (app *VideoSrt) Run(video string) {
	if video == "" {
		panic("enter a video file waiting to be processed .")
	}
	//1 校验视频
	if ValidVideo(video) != true {
		panic("the input video file does not exist .")
	}
	Log("提取音频文件 ...")
	//2 分离视频音频
	mp3path := ExtractVideoAudio(video, ".mp3")

	app.RunWithoutExtract(mp3path)

}

//直接翻译mp3
func (app *VideoSrt) RunWithoutExtract(mp3path string) {
	Log("上传音频文件 ...")
	//3.1 上传音频至OSS
	filelink := app.UploadAudioToCloud(mp3path)
	//3.2 获取完整链接
	singinurl, err := app.AliyunOss.GetObjectFileUrl(filelink)
	if err != nil {
		return
	}
	output := GetFileBaseName(mp3path)
	app.RunWithoutUpload(singinurl, output)
}

//不用上传mp3，有url直接翻译
func (app *VideoSrt) RunWithoutUpload(mp3url string, outputPath string) {
	Log("上传文件成功 , 识别中 ...")
	//4 阿里云录音文件识别
	AudioResult, err := app.AliyunAudioRecognition(mp3url)
	if err == nil {
		Log("文件识别成功 , 字幕处理中 ...")
		//5.输出字幕文件
		app.SaveResultToSrtFile(outputPath, AudioResult)
		Log("完成")
	} else {
		Log("文件识别失败")
	}
}

//提取视频音频文件
func ExtractVideoAudio(video string, audioType string) string {
	mp3path, err := ffmpeg.ExtractAudio(video, audioType)
	if err != nil {
		panic(err)
	}
	return mp3path
}

//上传音频至oss
func (app *VideoSrt) UploadAudioToCloud(audioFile string) string {
	target := app.AliyunOss
	bucketPath := app.AliyunOss.UploadPath + filepath.Base(audioFile)
	//上传
	if file, e := target.UploadFile(audioFile, bucketPath); e != nil {
		panic(e)
	} else {
		return file
	}
}

//阿里云录音文件识别
func (app *VideoSrt) AliyunAudioRecognition(filelink string) (AudioResult map[int64][]*cloud.AliyunAudioRecognitionResult, err error) {
	engine := app.AliyunCloud
	//创建识别请求
	taskid, e := engine.NewAudioFile(filelink)
	if e != nil {
		panic(e)
	}
	//查询请求
	return app.queryResultByTaskId(taskid)
}

func (app *VideoSrt) queryResultByTaskId(taskId string) (AudioResult map[int64][]*cloud.AliyunAudioRecognitionResult, err error) {
	AudioResult = make(map[int64][]*cloud.AliyunAudioRecognitionResult)
	intelligent_block := app.IntelligentBlock

	//遍历获取识别结果
	err = app.AliyunCloud.GetAudioFileResult(taskId, func(result []byte) {

		//结果处理
		statusText, _ := jsonparser.GetString(result, "StatusText") //结果状态
		if statusText == cloud.STATUS_SUCCESS {

			//智能分段
			if intelligent_block {
				cloud.AliyunAudioResultWordHandle(result, func(vresult *cloud.AliyunAudioRecognitionResult) {
					channelId := vresult.ChannelId

					_, isPresent := AudioResult[channelId]
					if isPresent {
						//追加
						AudioResult[channelId] = append(AudioResult[channelId], vresult)
					} else {
						//初始
						AudioResult[channelId] = []*cloud.AliyunAudioRecognitionResult{}
						AudioResult[channelId] = append(AudioResult[channelId], vresult)
					}
				})
				return
			}

			_, err := jsonparser.ArrayEach(result, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				text, _ := jsonparser.GetString(value, "Text")
				channelId, _ := jsonparser.GetInt(value, "ChannelId")
				beginTime, _ := jsonparser.GetInt(value, "BeginTime")
				endTime, _ := jsonparser.GetInt(value, "EndTime")
				silenceDuration, _ := jsonparser.GetInt(value, "SilenceDuration")
				speechRate, _ := jsonparser.GetInt(value, "SpeechRate")
				emotionValue, _ := jsonparser.GetInt(value, "EmotionValue")

				vresult := &cloud.AliyunAudioRecognitionResult{
					Text:            text,
					ChannelId:       channelId,
					BeginTime:       beginTime,
					EndTime:         endTime,
					SilenceDuration: silenceDuration,
					SpeechRate:      speechRate,
					EmotionValue:    emotionValue,
				}

				_, isPresent := AudioResult[channelId]
				if isPresent {
					//追加
					AudioResult[channelId] = append(AudioResult[channelId], vresult)
				} else {
					//初始
					AudioResult[channelId] = []*cloud.AliyunAudioRecognitionResult{}
					AudioResult[channelId] = append(AudioResult[channelId], vresult)
				}
			}, "Result", "Sentences")
			if err != nil {
				panic(err)
			}
		} else if statusText == cloud.STATUS_RUNNING {
			log.Print("云端执行中。。。")
		}
	})

	if err != nil {
		return nil, err
	}
	return AudioResult, nil
}

//阿里云录音识别结果集生成字幕文件
func (app *VideoSrt) SaveResultToSrtFile(srtPath string, AudioResult map[int64][]*cloud.AliyunAudioRecognitionResult) {

	for channel, result := range AudioResult {
		thisfile := srtPath + "_channel_" + strconv.FormatInt(channel, 10) + ".srt"
		//输出字幕文件
		println(thisfile)

		file, e := os.Create(thisfile)
		if e != nil {
			panic(e)
		}

		defer file.Close() //defer

		index := 0
		for _, data := range result {
			linestr := MakeSubtitleText(index, data.BeginTime, data.EndTime, data.Text)

			file.WriteString(linestr)

			index++
		}
	}
}

//拼接字幕字符串
func MakeSubtitleText(index int, startTime int64, endTime int64, text string) string {
	var content bytes.Buffer
	content.WriteString(strconv.Itoa(index))
	content.WriteString("\n")
	content.WriteString(SubtitleTimeMillisecond(startTime))
	content.WriteString(" --> ")
	content.WriteString(SubtitleTimeMillisecond(endTime))
	content.WriteString("\n")
	content.WriteString(text)
	content.WriteString("\n")
	content.WriteString("\n")
	return content.String()
}
