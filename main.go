package main

import (
	"log"
	"os"
	"taine-api/handler"
	"taine-api/infra"
	"taine-api/infra/postgres"
	"taine-api/interface/middleware"
	"taine-api/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	allowed := os.Getenv("ALLOWED_ORIGIN")
	if allowed == "" {
		allowed = "http://localhost:3000"
	}

	db, err := infra.NewDB(os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	// CORS（Nextのオリジンだけ許可）
	cfg := cors.Config{
		AllowOrigins: []string{allowed},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}
	router.Use(cors.New(cfg))

	// リポジトリの初期化
	userRepository := postgres.NewUserRepository(db.DB)
	orgRepository := postgres.NewOrganizationRepository(db.DB)
	membershipRepository := postgres.NewMembershipRepository(db.DB)
	tweetRepository := postgres.NewTweetRepository(db.DB)

	// サービスの初期化
	userService := usecase.NewUserService(userRepository)
	orgService := usecase.NewOrganizationSvc(orgRepository)
	membershipService := usecase.NewMembershipSvc(membershipRepository, userRepository, orgRepository)

	// ハンドラーの初期化
	webhookHandler := handler.NewWebhookHandler(userService, orgService, membershipService)
	router.POST("/webhooks/clerk", webhookHandler.Clerk)

	userUsecase := usecase.NewUserUsecase(userRepository)
	userHandler := handler.NewUserHandler(userUsecase)

	tweetUsecase := usecase.NewTweetUsecase(tweetRepository, userRepository)
	tweetHandler := handler.NewTweetHandler(tweetUsecase, userUsecase)

	router.GET("/api/health", func(c *gin.Context) { c.JSON(200, gin.H{"message": "OK"}) })

	// テスト用: 認証なしでアクセスできるエンドポイント
	router.GET("/api/test/tweets", tweetHandler.GetTweetsTest)

	api := router.Group("/api/v1", middleware.ClerkSessionAuth())
	api.GET("/me", userHandler.GetUserBySubID)

	// Tweet routes
	api.POST("/tweets", tweetHandler.CreateTweet)
	api.GET("/tweets", tweetHandler.GetTweets)      // 全てのtweetを取得（認証テスト用）
	api.GET("/tweets/my", tweetHandler.GetMyTweets) // 自分のtweetのみ取得
	api.GET("/tweets/:id", tweetHandler.GetTweetByID)
	api.PUT("/tweets/:id", tweetHandler.UpdateTweet)
	api.DELETE("/tweets/:id", tweetHandler.DeleteTweet)

	router.Run(":8080")
}
