package handlers

import (
	"RealWorldWeb/logger"
	"RealWorldWeb/middlewares"
	"RealWorldWeb/models"
	"RealWorldWeb/params/response"
	"RealWorldWeb/security"
	"RealWorldWeb/storage"
	"RealWorldWeb/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

func AddArticleCommentHandler(r *gin.Engine) {
	commentsGroup := r.Group("/api/article/:slug/comment")
	commentsGroup.GET("", GetArticleComment)

	commentsGroup.Use(middlewares.AuthMiddleware)
	commentsGroup.POST("", CreateArticleComment)
	commentsGroup.DELETE(":commentId", DeleteArticleComment)
}

func GetArticleComment(ctx *gin.Context) {
	log := logger.New(ctx)
	slug := ctx.Param("slug")
	log.Info("Get comment by slug: ", slug)

	articleComments, err := storage.GetArticleCommentByArticleId(ctx, slug)
	if err != nil {
		log.WithError(err).Infof("GetArticleCommentByArticleId failed, slug: %s", slug)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var res []*response.ArticleComment
	for _, articleComment := range articleComments {
		resArticleComment := response.ArticleComment{}
		resArticleComment.FromDB(&articleComment)
		res = append(res, &resArticleComment)
	}
	log.WithError(err).Infof("GetArticleCommentByArticleId res: %s", utils.JsonMarshal(res))

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"comments": res,
	})
}

func CreateArticleComment(ctx *gin.Context) {
	//req := make(map[string]interface{})
	log := logger.New(ctx)

	slug := ctx.Param("slug")
	article, err2 := storage.GetArticleBySlug(ctx, slug)
	if err2 != nil {
		log.WithError(err2).Infof("GetArticleBySlug failed, slug: %s", slug)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	req := struct {
		Comment struct {
			Body string `json:"body"`
		} `json:"comment"`
	}{}

	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		log.WithError(err).Infof("BindJSON failed, slug: %s", slug)
		return
	}

	creator := security.GetCurrentUsername(ctx)

	articleComment := &models.ArticleComment{
		AuthorUsername: creator,
		Body:           req.Comment.Body,
		ArticleId:      article.Id,
	}

	log.Infof("CreateArticleComment articleComment=%s", utils.JsonMarshal(articleComment))
	log.Infof("CreateArticleComment article=%s", utils.JsonMarshal(article))

	err2 = storage.CreateArticleComment(ctx, articleComment)
	if err2 != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		log.WithError(err).Infof("CreateArticleComment failed, slug: %s", slug)
		return
	}

	comment, err2 := storage.GetArticleCommentById(ctx, articleComment.Id)
	if err2 != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		log.WithError(err).Infof("GetArticleCommentById failed, slug: %s; id=%v", slug, articleComment.Id)
		return
	}

	res := response.ArticleComment{}
	res.FromDB(&comment)

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"comment": res,
	})
}

func DeleteArticleComment(ctx *gin.Context) {
	commentId := cast.ToInt64(ctx.Param("commentId"))
	log := logger.New(ctx)
	article, err := storage.GetArticleCommentById(ctx, commentId)
	if err != nil {

		if storage.IsNotFound(err) {
			ctx.AbortWithStatus(http.StatusNotFound)
			log.WithError(err).Infof("article is not fould failed, slug: %v", commentId)
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			log.WithError(err).Infof("article is not fould failed, slug: %v", commentId)
		}

		return
	}

	log.Infof("article=%s", utils.JsonMarshal(article))
	if article.AuthorUsername != security.GetCurrentUsername(ctx) {
		ctx.AbortWithStatus(http.StatusForbidden)
		log.WithError(err).Infof("AuthorUsername is not allowed, comment: %v; article.AuthorUsername=%s; currentUsername=%s", commentId, article.AuthorUsername, security.GetCurrentUsername(ctx))
		return
	}

	err = storage.DeleteArticleComment(ctx, commentId)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		log.WithError(err).Infof("DeleteArticleComment failed, slug: %v", commentId)
		return
	}

	ctx.Status(http.StatusOK)
}
