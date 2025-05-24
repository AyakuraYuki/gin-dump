package gin_dump

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func performRequest(r http.Handler, method, contentType, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,en-GB;q=0.6,en-US;q=0.5")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestMIMEJSON(t *testing.T) {
	router := gin.New()
	router.Use(DumpFunc(func(dumpStr string) {
		fmt.Println(dumpStr)
	}))

	router.POST("/dump", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"data": "gin-dump",
		})
	})

	type params struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	var payload = params{
		StartTime: "2025-05-24",
		EndTime:   "2025-05-24",
	}
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	body := bytes.NewBuffer(b)
	performRequest(router, "POST", gin.MIMEJSON, "/dump", body)
}

func TestMIMEPOSTFORM(t *testing.T) {
	router := gin.New()
	opts := []Option{
		WithShowRaw(true),
		WithCallback(func(dumpStr string) {
			fmt.Println(dumpStr)
		}),
	}
	router.Use(DumpWithOptions(opts...))

	router.POST("/dump", func(c *gin.Context) {
		bts, err := httputil.DumpRequest(c.Request, true)
		fmt.Println(string(bts), err)

		c.JSON(http.StatusOK, gin.H{
			"ok": true,
			"data": map[string]interface{}{
				"name": "jfise",
				"addr": "tpkeeper@qq.com",
			},
		})
	})

	form := make(url.Values)
	form.Set("foo", "bar")
	form.Add("foo", "bar2")
	form.Set("bar", "baz")

	body := strings.NewReader(form.Encode())
	performRequest(router, "POST", gin.MIMEPOSTForm, "/dump", body)
}
