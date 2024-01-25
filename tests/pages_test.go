package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomePage(t *testing.T) {
	baseURL := "http://localhost:3000"
	resp, err := http.Get(baseURL + "/")

	assert.NoError(t, err, "有错误发生， err 不为空")
	assert.Equal(t, resp.StatusCode, http.StatusOK, "应返回状态码 200")
}
