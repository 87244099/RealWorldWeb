package handlers

import (
	"RealWorldWeb/cache"
	"RealWorldWeb/logger"
	"RealWorldWeb/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddTagsHandler(r *gin.Engine) {

	tagGroup := r.Group("/api/tags")
	tagGroup.GET("", listPopularTags)

}

func listPopularTags(ctx *gin.Context) {
	log := logger.New(ctx)
	tags, err := cache.GetPopularTags(ctx)
	if err != nil {
		log.WithError(err).Infof("Failed to get popular tags from cache")
		return
	}
	if len(tags) != 0 {

		log.Infof("get popular tags from cache")

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"tags": tags,
		})

		return
	}
	tags, err = storage.ListPopularTags()
	log.Infof("get popular tags from DB")

	if err != nil {
		log.WithError(err).Infof("Failed to get popular tags from DB")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"tags": tags,
	})
}
