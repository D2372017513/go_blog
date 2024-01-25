package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllPage(t *testing.T) {
	baseURL := "http://localhost:3000"
	var tests = []struct {
		method   string
		url      string
		expected int
	}{
		{"GET", "/", 200},
		{"GET", "/about", 200},
		{"GET", "/notfound", 404},
		{"GET", "/articles", 200},
		{"GET", "/articles/create", 200},
		{"GET", "/articles/3", 200},
		{"GET", "/articles/3/edit", 200},
		{"POST", "/articles/3", 200},
		{"POST", "/articles", 200},
		{"POST", "/articles/1/delete", 404},
	}

	for _, test := range tests {
		t.Logf("当前测试链接：%s\n", test.url)

		var (
			resp *http.Response
			err  error
		)

		switch test.method {
		case "POST":
			data := make(map[string][]string)
			resp, err = http.PostForm(baseURL+test.url, data)
		default:
			resp, err = http.Get(baseURL + test.url)
		}

		assert.NoError(t, err, fmt.Sprintf("请求 %s 报错\n", test.url))
		assert.Equal(t, test.expected, resp.StatusCode, fmt.Sprintf("url %s 应返回状态码 %d \n", test.url, test.expected))
	}
}
