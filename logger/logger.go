package logger

import (
	"context"
	"github.com/sirupsen/logrus"
)

func New(ctx context.Context) *logrus.Entry {
	return logrus.New().WithContext(ctx)
}

/*

func init(){
//支持自定义logger的日志format
	logrus.SetFormatter(&logrus.JSONFormatter{

	})
}*/
