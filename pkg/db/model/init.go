package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jhw66/myvideo_lab4/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到控制台
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: cfg.Mysql.Dsn,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("open gorm failed:%w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB failed:%w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.Mysql.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Mysql.ConnMaxLifetime))
	sqlDB.SetMaxIdleConns(cfg.Mysql.MaxIdleConns)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql failed:%w", err)
	}

	if err := automigrate(db); err != nil {
		return nil, fmt.Errorf("autoimgrate failed:%w", err)
	}

	return db, nil
}
