package handlers

import (
	"RealWorldWeb/logger"
	"RealWorldWeb/params/response"
	"RealWorldWeb/storage"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func AddArticleHandler(r *gin.Engine) {
	articleGroup := r.Group("/api/article")
	articleGroup.GET("", listArticles)
}

func listArticles(ctx *gin.Context) {
	log := logger.New(ctx)
	limit := cast.ToInt(ctx.Query("limit")) //cast query string into int type
	offset := cast.ToInt(ctx.Query("offset"))
	log.Infof("limit:%d, offset:%d", limit, offset)
	articles, err := storage.ListArticles(ctx, limit, offset)
	if err != nil {
		log.Error(err)
	}
	total, err := storage.CountArticles(ctx)

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
}
