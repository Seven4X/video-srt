package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFileBaseName(t *testing.T) {
	println(GetFileBaseName("/Users/seven/Movies/ApowerREC/20201119_221121.mp3"))
}

func TestGetMp3FileList(t *testing.T) {
	res := GetMp3FileList("abc.txt")

	assert.NotNil(t, res)
}
