package storage

import (
	"RealWorldWeb/models"
	"context"
)

func CreateArticle(ctx context.Context, article *models.Article) error {
	return gormDB.WithContext(ctx).Create(article).Error
}
