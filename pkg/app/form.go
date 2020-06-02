package app

import (
	"github.com/coreos/pkg/capnslog"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
	e "immortality-demo/pkg/error"
	"runtime/debug"
)

var (
	ulog  = capnslog.NewPackageLogger("immortality", "app")
	vd    *validator.Validate
	trans ut.Translator
)

func AbortWithError(status string, err error) {
	panic(&Response{Status: status, Message: err.Error()})
}

func FaultWrap() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(*Response); ok {
					c.AbortWithStatusJSON(400,
						Response{Status: e.Status, Message: e.Message})
				} else {
					c.AbortWithStatusJSON(500,
						Response{Status: e.Status, Message: e.Message})
				}
			}
		}()
		c.Next()
	}
}

func ValidateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		zh := zh.New()
		trans, _ := ut.New(zh, zh).GetTranslator("zh")
		vd = validator.New()
		_ = zh_translations.RegisterDefaultTranslations(vd, trans)
		//自定义错误内容
		_ = vd.RegisterTranslation("required", trans, func(ut ut.Translator) error {
			return ut.Add("required", "{0} 字段不能为空!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		})
	}
}

// BindAndValid binds and validates data
func LoadBody(c *gin.Context, obj interface{}) {
	if err := c.Bind(obj); err != nil {
		ulog.Debugf("%s", debug.Stack())
		ulog.Error(err)
		AbortWithError(e.ParseBodyFailed, err)
	}

	if err := vd.Struct(obj); err != nil {
		err := err.(validator.ValidationErrors)
		comErr := NewValidatorError(err)
		ulog.Error(comErr)
		AbortWithError(e.ParseBodyFailed, err)
	}
}

func NewValidatorError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, e := range errs {
		res.Errors[e.Field()] = e.Translate(trans)
	}
	return res
}

func GetHeaderInfo(c *gin.Context) HeaderInfo {
	userId := c.GetHeader("X-User-Id")
	requestId := c.GetHeader("RequestId")
	if userId == "" || requestId == "" {
		msg := "Missing parameter UserId or RequestId"
		ulog.Error(msg)
		AbortWithError(e.InvalidParams, errors.New(msg))
	}

	headerInfo := HeaderInfo{
		UserId:    userId,
		RequestId: requestId,
	}
	return headerInfo
}
