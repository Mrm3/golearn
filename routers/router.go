package routers

import (
	"github.com/coreos/pkg/capnslog"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"immortality-demo/config"
	_ "immortality-demo/docs"
	"immortality-demo/pkg/app"
	v1 "immortality-demo/routers/api/v1"
)

var (
	ulog   = capnslog.NewPackageLogger("immortality", "api/v1")
	router *gin.Engine
)

func Router() *gin.Engine {
	if router == nil {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
		logfile := "/var/log/cvk-agent/compute.log"
		writer := config.GetWriter(logfile)
		loggerConfig := gin.LoggerConfig{Output: writer}
		router.Use(gin.LoggerWithConfig(loggerConfig), app.FaultWrap())
		router.Use(app.ValidateHandler())
		//router.Use(gin.Recovery())
	}
	return router
}

func init() {
	router := Router()

	//apiv1 := router.Group("/api/v1")
	////apiv1.Use(jwt.JWT())
	//{
	//	apiv1.GET("/tags", v1.GetTags)
	//	apiv1.POST("/tags", v1.AddTag)
	//	apiv1.PUT("/tags/:id", v1.EditTag)
	//	apiv1.DELETE("/tags/:id", v1.DeleteTag)
	//}
	test := router.Group("/v1")
	{
		test.Any("/test", v1.Test)
	}

	disk := router.Group("/v1")
	{
		disk.Any("/disks", v1.CreateDisk)
	}



	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
