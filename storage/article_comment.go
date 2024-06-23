package storage

import (
	"RealWorldWeb/models"
	"github.com/gin-gonic/gin"
)

func CreateArticleComment(ctx *gin.Context, articleComment *models.ArticleComment) error {
	return gormDB.WithContext(ctx).Create(articleComment).Error
}

func DeleteArticleComment(ctx *gin.Context, articleCommentId int64) error {
	return gormDB.WithContext(ctx).Delete(&models.ArticleComment{}, articleCommentId).Error
}

func GetArticleCommentByArticleId(ctx *gin.Context, articleSlug string) ([]models.ArticleComment, error) {
	comments := make([]models.ArticleComment, 0)
	err := gormDB.WithContext(ctx).
		Select("article_comment.*, user.email as author_user_email, user.bio as author_user_bio, user.image as author_user_image").
		Joins("LEFT JOIN user ON article_comment.author_username = user.username").
		Joins("INNER JOIN article ON article.id = article_comment.article_id").
		Where("article.slug = ?", articleSlug).
		Find(&comments).Error

	if err != nil {
		return nil, err
	}
	return comments, err
}

func GetArticleCommentById(ctx *gin.Context, commentId int64) (models.ArticleComment, error) {
	comment := models.ArticleComment{}
	err := gormDB.WithContext(ctx).
		Select("article_comment.*, user.email as author_user_email, user.bio as author_user_bio, user.image as author_user_image").
		Joins("LEFT JOIN user ON article_comment.author_username = user.username").
		Where("article_comment.id = ?", commentId).
		Find(&comment).Error
	if err != nil {
		return comment, err
	}
	return comment, nil
}
