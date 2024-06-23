package storage

import (
	"RealWorldWeb/models"
	"RealWorldWeb/params/request"
	"context"
	"gorm.io/gorm"
)

func CreateArticle(ctx context.Context, article *models.Article) error {
	return gormDB.WithContext(ctx).Create(article).Error
}

func ListArticles(ctx context.Context, req request.ListArticleQuery) ([]*models.Article, error) {
	var articles []*models.Article

	tx := listArticleTx(ctx, &req)

	err := tx.Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func CountArticles(ctx context.Context, req request.ListArticleQuery) (int64, error) {
	var count int64
	//gormDb.Table("article").Count(&count);
	err := listArticleTx(ctx, &req).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func listArticleTx(ctx context.Context, req *request.ListArticleQuery) *gorm.DB {

	tx := gormDB.WithContext(ctx).Model(models.Article{}).
		Select("article.*, user.email as author_user_email, user.bio as author_user_bio, user.image as author_user_image").
		Joins("LEFT JOIN user ON article.author_username = user.username").
		Order("created_at desc").
		Offset(req.Offset).Limit(req.Limit)

	// like会有性能问题，但是先这样子，因为这里的tag字段属于数组json，只能扫全表
	if req.Tag != "" {
		tx = tx.Where("article.tag_list like ?", "%"+req.Tag+"%")
	}
	return tx
}

func GetArticleBySlug(ctx context.Context, slug string) (*models.Article, error) {
	var article models.Article
	err := gormDB.WithContext(ctx).Model(models.Article{}).
		Select("article.*, user.email as author_user_email, user.bio as author_user_bio, user.image as author_user_image").
		Joins("LEFT JOIN user ON article.author_username = user.username").
		Where("slug = ?", slug).
		First(&article).Error
	return &article, err
}

func UpdateArticle(ctx context.Context, slug string, article *models.Article) error {
	return gormDB.WithContext(ctx).Where("slug=?", slug).Updates(article).Error
}

func DeleteArticle(ctx context.Context, slug string) error {
	//why pass it ?
	return gormDB.WithContext(ctx).Where("slug=?", slug).Delete(&models.Article{}).Error
}
