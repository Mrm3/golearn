package v1

import (
	"github.com/coreos/pkg/capnslog"
	"github.com/gin-gonic/gin"
	"immortality-demo/pkg/app"
)

var ulog = capnslog.NewPackageLogger("immortality", "api/v1/test")

// @Summary Get Auth
// @Tags Test
// @Produce  json
// @Param body body common.Response true "test api"
// @Success 200 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /test [get]
func Test(c *gin.Context) {
	var respData app.Response
	ulog.Println("this is api test")
	//参数校验
	if err := c.Bind(&respData); err != nil {
		ulog.Println(err.Error())
		return
	}
	c.JSON(200, gin.H{
		"success": respData.Status,
		"msg":     respData.Message,
		"data":    respData.Detail,
	})
}
