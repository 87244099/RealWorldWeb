package handlers

import (
	"RealWorldWeb/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddTagsHandler(r *gin.Engine) {

	tagGroup := r.Group("/api/tags")
	tagGroup.GET("", listPopularTags)

}

func listPopularTags(ctx *gin.Context) {
	tags, err := storage.ListPopularTags()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"tags": tags,
	})
}
