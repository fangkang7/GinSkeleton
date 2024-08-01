package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"goskeleton/app/global/variable"
	"goskeleton/app/utils/perror"
)

type LogData struct {
	logType   zapcore.Level
	platform  string
	url       string
	method    string
	params    interface{}
	companyId int64
	userId    int64
	code      int
	time      int64
	msg       string
	data      interface{}
	traceId   int64
}

func writeLog(l LogData) {
	jsonData, _ := json.Marshal(l.data)
	jsonParams, _ := json.Marshal(l.params)
	variable.ApiLog.Log(
		l.logType,
		l.msg,
		zap.Int("code", l.code),
		zap.String("platform", l.platform),
		zap.Int64("companyId", l.companyId),
		zap.Int64("userId", l.userId),
		zap.String("url", l.url),
		zap.String("method", l.method),
		zap.ByteString("params", jsonParams),
		zap.ByteString("data", jsonData),
		zap.Int64("traceId", l.traceId),
		zap.Int64("time", l.time),
	)

}

func Info(con *gin.Context, msg string, data interface{}) {
	logInfo := LogData{
		logType:  zap.InfoLevel,
		platform: con.GetString("platform"),
		msg:      msg,
		data:     data,
		traceId:  con.GetInt64("traceId"),
	}
	writeLog(logInfo)
}

func Error(con *gin.Context, err perror.Error) {
	logInfo := LogData{
		logType:  zap.ErrorLevel,
		platform: con.GetString("platform"),
		msg:      err.Error(),
		traceId:  con.GetInt64("traceId"),
	}
	writeLog(logInfo)
}

/**
*context.Reauest.Body 获取一次后会被清空
*RequestLog 在请求时打印params，返回时不打印
 */
func RequestStartLog(con *gin.Context, msg string, data interface{}) {
	logInfo := LogData{
		logType:   zap.InfoLevel,
		platform:  con.GetString("platform"),
		url:       con.Request.URL.Path,
		method:    con.Request.Method,
		params:    getQueryParams(con),
		msg:       msg,
		data:      data,
		companyId: con.GetInt64("companyId"),
		userId:    con.GetInt64("userId"),
		time:      (time.Now().UnixNano() / 1e6) - con.GetInt64("startTime"),
		traceId:   con.GetInt64("traceId"),
	}
	writeLog(logInfo)
}

func RequestEndLog(con *gin.Context, code int, msg string, data interface{}) {
	logInfo := LogData{
		logType:   zap.InfoLevel,
		platform:  con.GetString("platform"),
		url:       con.Request.URL.Path,
		method:    con.Request.Method,
		code:      code,
		msg:       msg,
		data:      data,
		companyId: con.GetInt64("companyId"),
		userId:    con.GetInt64("userId"),
		time:      (time.Now().UnixNano() / 1e6) - con.GetInt64("startTime"),
		traceId:   con.GetInt64("traceId"),
	}
	writeLog(logInfo)
}

func getQueryParams(c *gin.Context) map[string]interface{} {
	queryMap := map[string]interface{}{
		"GET":  getParamsGet(c),
		"POST": getParamsPost(c),
		"json": getParamsJson(c),
	}
	return queryMap
}
func getParamsGet(c *gin.Context) map[string]interface{} {
	query := c.Request.URL.Query()
	var queryMap = make(map[string]interface{}, len(query))
	for k := range query {
		queryMap[k] = c.Query(k)
	}
	return queryMap
}

func getParamsPost(c *gin.Context) map[string]interface{} {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		if !errors.Is(err, http.ErrNotMultipart) {
			return map[string]interface{}{
				"err": "参数获取失败",
			}
		}
	}
	var postMap = make(map[string]interface{}, len(c.Request.PostForm))
	for k, v := range c.Request.PostForm {
		if len(v) > 1 {
			postMap[k] = v
		} else if len(v) == 1 {
			postMap[k] = v[0]
		}
	}

	return postMap
}

// json传参
func getParamsJson(c *gin.Context) map[string]interface{} {
	reqBytes, _ := c.GetRawData() // 请求包体写回。
	if len(reqBytes) > 0 {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBytes))
	}
	var body = make(map[string]interface{})
	json.Unmarshal(reqBytes, &body)
	return body
}
