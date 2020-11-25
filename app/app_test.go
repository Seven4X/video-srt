package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery(t *testing.T) {

	app := NewByConfig("../config.ini")

	res, err := app.queryResultByTaskId("6687764a2ef911eb9a720163bb30af52")
	app.SaveResultToSrtFile("/Users/seven/data/demo", res)
	assert.Nil(t, err)
	assert.NotNil(t, res)

}
