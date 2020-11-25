package ffmpeg

import (
	"path/filepath"
	"testing"
)

func Test_ExtractAudio(t *testing.T) {

	ExtractAudio("/Users/seven/Movies/ApowerREC/20201119_212227.mp4", ".mp3")

}

func TestPath(t *testing.T) {
	println(filepath.Base("/Users/seven/Movies/ApowerREC/20201119_212227.mp4"))
	println(filepath.Ext("/Users/seven/Movies/ApowerREC/20201119_212227.mp4"))
	println(filepath.Dir("/Users/seven/Movies/ApowerREC/20201119_212227.mp4"))
}

func TestGetExecPath(t *testing.T) {
	getExecPath()
}
