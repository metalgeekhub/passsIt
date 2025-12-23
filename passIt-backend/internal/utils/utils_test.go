package utils

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestDecodeServerInput_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := TestStruct{Name: "test", Value: 123}
	jsonData, _ := json.Marshal(testData)

	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	var result TestStruct
	success := DecodeServerInput(c, &result)

	assert.True(t, success)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 123, result.Value)
}

func TestDecodeServerInput_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	var result TestStruct
	success := DecodeServerInput(c, &result)

	assert.False(t, success)
	assert.Equal(t, 400, w.Code)
}

func TestInsecureHttpContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ctx := InsecureHttpContext(c)
	
	assert.NotNil(t, ctx)
}
