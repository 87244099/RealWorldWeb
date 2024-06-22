package handlers

import (
	"RealWorldWeb/logger"
	"RealWorldWeb/params/request"
	"RealWorldWeb/params/response"
	"RealWorldWeb/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddArticleHandler(r *gin.Engine) {
	articleGroup := r.Group("/api/article")
	articleGroup.GET("", listArticles)
}

func listArticles(ctx *gin.Context) {
	log := logger.New(ctx)

	//读取特定的queryString到req上
	var req request.ListArticleQuery
	err2 := ctx.BindQuery(&req)
	if err2 != nil {
		ctx.AbortWithError(http.StatusBadRequest, err2)
		return
	}

	log.Infof("limit:%d, offset:%d, tag=%s", req.Limit, req.Offset, req.Tag)

	articles, err := storage.ListArticles(ctx, req)
	if err != nil {
		log.Error(err)
	}
	total, err := storage.CountArticles(ctx, req)

	var res response.ListArticlesResponse
	res.ArticlesCount = total
	for _, article := range articles {
		res.Articles = append(res.Articles, &response.Article{
			Author: &response.ArticleAuthor{
				Bio:       article.AuthorUserBio,
				Following: false,
				Image:     article.AuthorUserImage,
				Username:  article.AuthorUsername,
			},
			Title:          article.Title,
			Slug:           article.Slug,
			Body:           article.Body,
			Description:    article.Description,
			TagList:        article.TagList,
			Favorited:      false,
			FavoritesCount: 0,
			CreatedAt:      article.CreatedAt,
			UpdatedAt:      article.UpdatedAt,
		})
	}
	ctx.JSON(http.StatusOK, res)
}
