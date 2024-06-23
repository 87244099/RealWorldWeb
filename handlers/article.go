package handlers

import (
	"RealWorldWeb/logger"
	"RealWorldWeb/middlewares"
	"RealWorldWeb/models"
	"RealWorldWeb/params/request"
	"RealWorldWeb/params/response"
	"RealWorldWeb/security"
	"RealWorldWeb/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func AddArticleHandler(r *gin.Engine) {
	articleGroup := r.Group("/api/articles")
	articleGroup.GET("", listArticles)
	articleGroup.GET("/:slug", getArticle)

	articleGroup.Use(middlewares.AuthMiddleware)
	articleGroup.POST("", createArticle)
	articleGroup.PUT("/:slug", editArticle)
	articleGroup.DELETE("/:slug", deleteArticle)

}

func listArticles(ctx *gin.Context) {
	log := logger.New(ctx)

	//读取特定的queryString到req上
	var req request.ListArticleQuery
	err2 := ctx.BindQuery(&req)
	if err2 != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
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

		resArticle := &response.Article{}
		resArticle.FromDB(article)
		res.Articles = append(res.Articles, resArticle)
	}
	ctx.JSON(http.StatusOK, res)
}

func getArticle(ctx *gin.Context) {

	slug := ctx.Param("slug")
	article, err := storage.GetArticleBySlug(ctx, slug)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	resArticle := &response.Article{}
	resArticle.FromDB(article)
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"article": resArticle,
	})
}

func createArticle(ctx *gin.Context) {
	log := logger.New(ctx)
	var req request.CreateArticleRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		log.WithError(err).Infof("parse query failed")
		//time="2024-06-23T13:19:38+08:00" level=info msg="parse query failed" error="invalid character '-' in numeric literal"
		//生成的第一篇文章slug为：testTitle-8f06022b-9d87-41fc-92da-79d6c0008afb
		//postman 那边要选择raw格式发送
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	slug := strings.ReplaceAll(req.Title, " ", "-") + "-" + uuid.NewString()
	err = storage.CreateArticle(ctx, &models.Article{
		AuthorUsername: security.GetCurrentUsername(ctx),
		Title:          req.Title,
		Slug:           slug,
		Body:           req.Body,
		Description:    req.Description,
		TagList:        req.TagList,
	})

	if err != nil {
		log.WithError(err).Infof("create article failed")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	article, err := storage.GetArticleBySlug(ctx, slug)
	if err != nil {
		log.WithError(err).Infof("get article failed")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resArticle := &response.Article{}
	resArticle.FromDB(article)
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"article": resArticle,
	})

	//storage.CreateArticle()
}

func editArticle(ctx *gin.Context) {
	oldSlug := ctx.Param("slug")
	log := logger.New(ctx)
	var req request.CreateArticleRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		log.WithError(err).Infof("parse query failed")
		return
	}

	httpCode := checkPermissionWithSlug(ctx, oldSlug)
	if httpCode != http.StatusOK {
		ctx.AbortWithStatus(httpCode)
		log.Infof("has no permission; httpCode: %d", httpCode)
		return
	}

	slug := strings.ReplaceAll(req.Title, " ", "-") + "-" + uuid.NewString()
	err = storage.UpdateArticle(ctx, oldSlug, &models.Article{
		AuthorUsername: security.GetCurrentUsername(ctx),
		Title:          req.Title,
		Slug:           slug,
		Body:           req.Body,
		Description:    req.Description,
		TagList:        req.TagList,
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		log.WithError(err).Infof("update article failed")
		return
	}

	article, err := storage.GetArticleBySlug(ctx, slug)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		log.WithError(err).Infof("get article failed")
		return
	}

	resArticle := &response.Article{}
	resArticle.FromDB(article)
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"article": resArticle,
	})

}

func deleteArticle(ctx *gin.Context) {
	slug := ctx.Param("slug")

	httpCode := checkPermissionWithSlug(ctx, slug)
	if httpCode != http.StatusOK {
		ctx.AbortWithStatus(httpCode)
		return
	}

	err := storage.DeleteArticle(ctx, slug)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func checkPermissionWithSlug(ctx *gin.Context, oldSlug string) int {
	oldArticle, err := storage.GetArticleBySlug(ctx, oldSlug)
	if err != nil {
		return http.StatusInternalServerError
	}
	if oldArticle.AuthorUsername != security.GetCurrentUsername(ctx) {
		return http.StatusForbidden
	}
	return http.StatusOK
}
