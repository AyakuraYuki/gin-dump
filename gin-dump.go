package gin_dump

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func Dump() gin.HandlerFunc {
	return dump(&options{
		showReq:     defaultShowReq,
		showRsp:     defaultShowRsp,
		showBody:    defaultShowBody,
		showHeaders: defaultShowHeaders,
		showCookies: defaultShowCookies,
		showRaw:     defaultShowRaw,
		cb:          nil,
	})
}

func DumpFunc(cb Callback) gin.HandlerFunc {
	return dump(&options{
		showReq:     defaultShowReq,
		showRsp:     defaultShowRsp,
		showBody:    defaultShowBody,
		showHeaders: defaultShowHeaders,
		showCookies: defaultShowCookies,
		showRaw:     defaultShowRaw,
		cb:          cb,
	})
}

func DumpWithOptions(opts ...Option) gin.HandlerFunc {
	option := &options{
		showReq:     defaultShowReq,
		showRsp:     defaultShowRsp,
		showBody:    defaultShowBody,
		showHeaders: defaultShowHeaders,
		showCookies: defaultShowCookies,
		showRaw:     defaultShowRaw,
		cb:          nil,
	}
	for _, opt := range opts {
		opt(option)
	}
	return dump(option)
}

func dump(option *options) gin.HandlerFunc {
	hiddenHeaders := make([]string, 0)
	hiddenBody := make([]string, 0)

	if !option.showCookies {
		hiddenHeaders = append(hiddenHeaders, "cookie")
	}

	return func(ctx *gin.Context) {
		var strB strings.Builder

		if option.showReq && option.showHeaders {
			// dump request header
			s, err := FormatToJson(ctx.Request.Header, hiddenHeaders, true)
			if err != nil {
				strB.WriteString(fmt.Sprintf("\nparse request header error: %s\n", err.Error()))
			} else {
				strB.WriteString("Request-Header:\n")
				strB.WriteString(string(s))
			}
			strB.WriteString("\n")
		}

		if option.showReq && option.showBody {
			// dump request body
			if ctx.Request.ContentLength > 0 {
				buf, err := io.ReadAll(ctx.Request.Body)
				if err != nil {
					strB.WriteString(fmt.Sprintf("\nread bodyCache err: %s\n", err.Error()))
					goto DumpRes
				}

				rdr := io.NopCloser(bytes.NewBuffer(buf))
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
				ctGet := ctx.Request.Header.Get("Content-Type")
				ct, _, err := mime.ParseMediaType(ctGet)
				if err != nil {
					strB.WriteString(fmt.Sprintf("\ncontent_type: %s, parse err: %s\n", ctGet, err.Error()))
					goto DumpRes
				}

				switch ct {
				case gin.MIMEJSON:
					bts, err := io.ReadAll(rdr)
					if err != nil {
						strB.WriteString(fmt.Sprintf("\nread rdr err: %s\n", err.Error()))
						goto DumpRes
					}

					if option.showRaw {
						strB.WriteString("\nRequest-Body (raw):\n")
						strB.WriteString(string(bts) + "\n")
					}

					s, err := FormatJsonBytes(bts, hiddenBody, false)
					if err != nil {
						strB.WriteString(fmt.Sprintf("\nparse request body err: %s\n", err.Error()))
						goto DumpRes
					}

					strB.WriteString("\nRequest-Body:\n")
					strB.WriteString(string(s) + "\n")

				case gin.MIMEPOSTForm:
					bts, err := io.ReadAll(rdr)
					if err != nil {
						strB.WriteString(fmt.Sprintf("\nread rdr err: %s\n", err.Error()))
						goto DumpRes
					}

					if option.showRaw {
						strB.WriteString("\nRequest-Body (raw):\n")
						strB.WriteString(string(bts) + "\n")
					}

					val, err := url.ParseQuery(string(bts))
					s, err := FormatToJson(val, hiddenBody, false)
					if err != nil {
						strB.WriteString(fmt.Sprintf("\nparse request body err: %s\n", err.Error()))
						goto DumpRes
					}
					strB.WriteString("\nRequest-Body:\n")
					strB.WriteString(string(s) + "\n")

				case gin.MIMEMultipartPOSTForm:
				default:
				}
			}

		DumpRes:
			ctx.Writer = &bodyWriter{bodyCache: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
			ctx.Next()
		}

		if option.showRsp && option.showHeaders {
			// dump response header
			s, err := FormatToJson(ctx.Writer.Header(), hiddenHeaders, true)
			if err != nil {
				strB.WriteString(fmt.Sprintf("\nparse response header err: %s\n", err.Error()))
			} else {
				strB.WriteString("\nResponse-Header:\n")
				strB.WriteString(string(s))
			}
			strB.WriteString("\n")
		}

		if option.showRsp && option.showBody {
			bw, ok := ctx.Writer.(*bodyWriter)
			if !ok {
				strB.WriteString("\nbodyWriter was override, can not read bodyCache")
				goto End
			}

			// dump response body
			if bodyAllowedForStatus(ctx.Writer.Status()) && bw.bodyCache.Len() > 0 {
				ctGet := ctx.Writer.Header().Get("Content-Type")
				ct, _, err := mime.ParseMediaType(ctGet)
				if err != nil {
					strB.WriteString(fmt.Sprintf("\ncontent-type: %s parse, err \n %s", ctGet, err.Error()))
					goto End
				}
				switch ct {
				case gin.MIMEJSON:
					cachedBs := bw.bodyCache.Bytes()
					if option.showRaw {
						strB.WriteString("\nResponse-Body (raw):\n")
						strB.WriteString(string(cachedBs) + "\n")
					}
					s, err := FormatJsonBytes(cachedBs, hiddenBody, false)
					if err != nil {
						strB.WriteString(fmt.Sprintf("\nparse bodyCache err: %s\n", err.Error()))
						goto End
					}
					strB.WriteString("\nResponse-Body:\n")
					strB.WriteString(string(s) + "\n")

				case gin.MIMEHTML:
				default:
				}
			}
		}

	End:
		if option.cb != nil {
			option.cb(strB.String())
		} else {
			fmt.Print(strB.String())
		}
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	bodyCache *bytes.Buffer
}

// override Write()
func (w bodyWriter) Write(b []byte) (int, error) {
	w.bodyCache.Write(b)
	return w.ResponseWriter.Write(b)
}

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
