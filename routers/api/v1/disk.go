package v1

import (
	"github.com/gin-gonic/gin"
	"immortality-demo/pkg/app"
	e "immortality-demo/pkg/error"
	"immortality-demo/service/disk_service"
	"immortality-demo/service/model"
	"net/http"
)

// @Summary Create a disk
// @Description create server
// @Tags Servers
// @Accept  json
// @Produce  json
// @Param X-User-Id header string true "X-User-Id"
// @Param RequestId header string true "RequestId"
// @Param body body model.CreateRequest true "create server body"
// @Success 200 {object} common.Response
// @Router /v1/disk [post]
func CreateDisk(c *gin.Context) {
	var request model.CreateDiskParams
	appG := app.Gin{C: c}

	app.LoadBody(c, &request)
	header := app.GetHeaderInfo(c)

	result, err := (disk_service.CreateDiskHandler{}).Handle(&request, header.UserId, header.RequestId)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.SUCCESS, "", result)
	}
	appG.Response(http.StatusOK, e.CreateDiskError, "", result)
}
