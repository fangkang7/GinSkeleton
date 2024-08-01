package response

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/perror"
	"goskeleton/app/utils/tool"
	"goskeleton/app/utils/validator_translation"
	"net/http"
	//"goskeleton/app/utils/logger"
	"strings"
)

type SuccessData struct {
	Code int
	Msg  string
	Data interface{}
}

func ReturnData(code int, msg string, data interface{}) (*SuccessData, perror.Error) {
	return &SuccessData{
		Code: code,
		Msg:  msg,
		Data: data,
	}, nil
}

func ReturnSuccessData(data interface{}) (*SuccessData, perror.Error) {
	return &SuccessData{
		Code: consts.CurdStatusOkCode,
		Msg:  consts.CurdStatusOkMsg,
		Data: data,
	}, nil
}

func ReturnError(err perror.Error) (*SuccessData, perror.Error) {
	return &SuccessData{
		Code: -1,
		Msg:  "请求失败",
		Data: nil,
	}, err
}

func returnJson(Context *gin.Context, httpCode int, dataCode int, msg string, data interface{}) {

	//Context.Header("key2020","value2020")  	//可以根据实际情况在头部添加额外的其他信息
	Context.JSON(httpCode, gin.H{
		"code":    dataCode,
		"msg":     msg,
		"data":    data,
		"traceId": Context.GetInt64("traceId"),
		"env":     tool.StringBuild(variable.ConfigYml.GetString("APPNAME"), "-", variable.ConfigYml.GetString("APPENV")),
	})
}

// start 业务一般只使用RequestSuccess 和 RequestFail 方法，因为我们的业务请求的http状态200时前端才会处理
// RequestSuccess 直接返回成功
func RequestSuccess(c *gin.Context, d *SuccessData) {
	returnJson(c, http.StatusOK, d.Code, d.Msg, d.Data)
	if d.Code != consts.AuthorizationFaildCode {
		//logger.RequestEndLog(c, d.Code, d.Msg, d.Data)
	}
}

// RequestFail 报错直接返回
func RequestFail(c *gin.Context, e perror.Error) {
	appDebug := variable.ConfigYml.GetBool("AppDebug")
	var errMsg string
	if appDebug {
		errMsg = e.Error()
	} else {
		errMsg = e.Msg()
	}
	if errMsg == "" {
		errMsg = "系统错误"
	}
	returnJson(c, http.StatusOK, e.Code(), errMsg, nil)
	if e.Code() != consts.AuthorizationFaildCode {
		//logger.RequestEndLog(c, e.Code(), e.Error(), nil)
	}
	c.Abort()
}

func ReturnJson(Context *gin.Context, httpCode int, dataCode int, msg string, data interface{}) {

	//Context.Header("key2020","value2020")  	//可以根据实际情况在头部添加额外的其他信息
	Context.JSON(httpCode, gin.H{
		"code": dataCode,
		"msg":  msg,
		"data": data,
	})
}

// ReturnJsonFromString 将json字符窜以标准json格式返回（例如，从redis读取json格式的字符串，返回给浏览器json格式）
func ReturnJsonFromString(Context *gin.Context, httpCode int, jsonStr string) {
	Context.Header("Content-Type", "application/json; charset=utf-8")
	Context.String(httpCode, jsonStr)
}

// 语法糖函数封装

// Success 直接返回成功
func Success(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusOK, consts.CurdStatusOkCode, msg, data)
}

// Fail 失败的业务逻辑
func Fail(c *gin.Context, dataCode int, msg string, data interface{}) {
	ReturnJson(c, http.StatusBadRequest, dataCode, msg, data)
	c.Abort()
}

// ErrorTokenBaseInfo token 基本的格式错误
func ErrorTokenBaseInfo(c *gin.Context) {
	ReturnJson(c, http.StatusBadRequest, http.StatusBadRequest, my_errors.ErrorsTokenBaseInfo, "")
	//终止可能已经被加载的其他回调函数的执行
	c.Abort()
}

// ErrorTokenAuthFail token 权限校验失败
func ErrorTokenAuthFail(c *gin.Context) {
	ReturnJson(c, http.StatusUnauthorized, http.StatusUnauthorized, my_errors.ErrorsNoAuthorization, "")
	//终止可能已经被加载的其他回调函数的执行
	c.Abort()
}

// ErrorTokenRefreshFail token不符合刷新条件
func ErrorTokenRefreshFail(c *gin.Context) {
	ReturnJson(c, http.StatusUnauthorized, http.StatusUnauthorized, my_errors.ErrorsRefreshTokenFail, "")
	//终止可能已经被加载的其他回调函数的执行
	c.Abort()
}

// token 参数校验错误
func TokenErrorParam(c *gin.Context, wrongParam interface{}) {
	ReturnJson(c, http.StatusUnauthorized, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, wrongParam)
	c.Abort()
}

// ErrorCasbinAuthFail 鉴权失败，返回 405 方法不允许访问
func ErrorCasbinAuthFail(c *gin.Context, msg interface{}) {
	ReturnJson(c, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed, my_errors.ErrorsCasbinNoAuthorization, msg)
	c.Abort()
}

// ErrorParam 参数校验错误
func ErrorParam(c *gin.Context, wrongParam interface{}) {
	ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, wrongParam)
	c.Abort()
}

// ErrorSystem 系统执行代码错误
func ErrorSystem(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusInternalServerError, consts.ServerOccurredErrorCode, consts.ServerOccurredErrorMsg+msg, data)
	c.Abort()
}

// ValidatorError 翻译表单参数验证器出现的校验错误
func ValidatorError(c *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		wrongParam := validator_translation.RemoveTopStruct(errs.Translate(validator_translation.Trans))
		ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, wrongParam)
	} else {
		errStr := err.Error()
		// multipart:nextpart:eof 错误表示验证器需要一些参数，但是调用者没有提交任何参数
		if strings.ReplaceAll(strings.ToLower(errStr), " ", "") == "multipart:nextpart:eof" {
			ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, gin.H{"tips": my_errors.ErrorNotAllParamsIsBlank})
		} else {
			ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, gin.H{"tips": errStr})
		}
	}
	c.Abort()
}
