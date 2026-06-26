package main

import (
	"fmt"
	"log"
	"time"

	"github.com/subosito/gotenv"

	"blog/api/middleware"
	"blog/api/router"
	"blog/config"
	"blog/internal/service"
	"blog/pkg/cache"
	"blog/pkg/database"
	"blog/pkg/logger"
	"blog/pkg/email"
	"blog/pkg/utils"

	_ "blog/docs" // swagger docs
)

// @title           Blog
// @version         v1
// @description     一个基于 Go 语言的现代化博客后端服务

// @contact.name   mccsjs
// @contact.email  mccsjs@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 在请求头中添加 Bearer Token，格式：Bearer {token}

func main() {
	// 加载 .env 文件
	_ = gotenv.Load()

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	db, err := database.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		_ = db.Close()
	}()
	defer logger.Close()
	defer middleware.ClosePanicLogFile()

	// 初始化 Redis 缓存
	var cacheClient *cache.Client
	if cfg.Redis.Addr != "" && cfg.Redis.Password != "" {
		cc, err := cache.NewClient(cache.Config{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
			TTL:      5 * time.Minute,
		})
		if err != nil {
			logger.Warn("Redis 连接失败，将不使用缓存: %v", err)
		} else {
			cacheClient = cc
			defer cc.Close()
			logger.Info("Redis 缓存初始化成功 (%s)", cfg.Redis.Addr)
		}
	} else {
		logger.Info("Redis 未配置，跳过缓存初始化")
	}

	// 执行数据库迁移
	if err := database.RunMigrations(db.DB); err != nil {
		middleware.ClosePanicLogFile()
		logger.Close()
		_ = db.Close()
		log.Fatalf("Failed to run migrations: %v", err) //nolint:gocritic // 已手动关闭资源
	}

	// 从数据库加载运行时配置（邮箱、上传等）
	settingService := service.NewSettingService(db.DB)
	settingService.SetConfig(cfg)               // 设置全局配置对象，用于热重载
	_ = settingService.ApplyDatabaseConfig(cfg) // 应用数据库配置

	// 初始化 IP 地理位置
	_ = utils.InitIPSearcher("")
	defer utils.CloseIPSearcher()


	// 初始化邮件服务
	emailClient := email.Initialize(cfg)
	if emailClient == nil {
		logger.Warn("邮件服务初始化失败，将无法发送邮件通知")
	} else {
		logger.Info("邮件服务初始化成功")
	}

	// 初始化路由
	r := router.InitRouter(db, cfg, cacheClient)

	// 启动服务器
	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port)
	logger.Info("Server is running at http://localhost:%d", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
