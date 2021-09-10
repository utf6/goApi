package main

import (
	"github.com/robfig/cron"
	"github.com/utf6/goApi/app/models"
	"github.com/utf6/goApi/pkg/logger"
	"time"
)

func init() {
	logger.Info("staring....")

	c := cron.New()
	c.AddFunc("* * * * * *", func() {
		logger.Info("执行清理标签数据")
		models.CleanTag()
	})

	c.AddFunc("* * * * * *", func() {
		logger.Info("执行清理文章数据")
		models.CleanArticle()
	})
	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
