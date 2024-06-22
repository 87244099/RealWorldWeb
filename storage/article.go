package storage

import (
	"RealWorldWeb/models"
	"context"
)

func CreateArticle(ctx context.Context, article *models.Article) error {
	return gormDB.WithContext(ctx).Create(article).Error
}

func ListArticles(ctx context.Context, limit, offset int) ([]*models.Article, error) {
	var articles []*models.Article
	err := gormDB.WithContext(ctx).Model(models.Article{}).
		Select("article.*, user.email as author_user_name, user.bio as author_user_bio, user.image as author_user_image").
		Joins("LEFT JOIN user on article.author_username = user.username").
		Order("created_at desc").Limit(limit).Offset(offset).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func CountArticles(ctx context.Context) (int64, error) {
	var count int64
	//gormDb.Table("article").Count(&count);
	err := gormDB.WithContext(ctx).Model(models.Article{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
