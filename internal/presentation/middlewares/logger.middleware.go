package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		var reqBody []byte
		var err error
		if c.Request.Body != nil {
			reqBody, err = io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		w := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		elapsed := time.Since(startTime)

		respBody := w.body.String()

		requestHeaderMap := map[string]string{}

		requestMap := map[string]interface{}{
			"url":     c.Request.URL.Path,
			"method":  c.Request.Method,
			"body":    string(reqBody),
			"headers": requestHeaderMap,
		}
		if q := c.Request.URL.RawQuery; q != "" {
			requestMap["query_string"] = q
		}

		responseMap := map[string]interface{}{
			"status_code": c.Writer.Status(),
			"body":        respBody,
		}

		log.Info().
			Int64("duration_ms", elapsed.Milliseconds()).
			Interface("request", requestMap).
			Interface("response", responseMap).
			Msgf("[%s] API Request and Response", viper.GetString("APP_NAME"))
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
