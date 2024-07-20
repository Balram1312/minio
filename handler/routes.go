package handler

import "github.com/gin-gonic/gin"

func (app *MyApp) Routes(r *gin.Engine) {
	r.GET("file", app.DownloadFile)
	r.GET("buckets", app.ListBuckets)
	r.POST("file", app.UploadFile)
	r.GET("objects", app.ListObjects)
}
