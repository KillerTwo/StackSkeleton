package web

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/web"
)

type UploadRouter struct {
}

func (s *UploadRouter) InitUploadRouter(router *gin.RouterGroup) {
	upload := router.Group("/")
	{
		upload.POST("upload/file", (&web.Upload{}).StartUpload)
	}
}
