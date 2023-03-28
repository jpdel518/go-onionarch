package main

import (
	"github.com/joho/godotenv"
	"github.com/jpdel518/go-onionarch/domain/service"
	"github.com/jpdel518/go-onionarch/infrastructure/rdb"
	"github.com/jpdel518/go-onionarch/infrastructure/rdb/mysql"
	handler2 "github.com/jpdel518/go-onionarch/presenter/handler"
	"github.com/jpdel518/go-onionarch/usecase"
	"github.com/jpdel518/go-onionarch/utils"
	"log"
	"os"
	"time"
)

func init() {
	loadEnv()
	initLogging()
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}
}

func initLogging() {
	logfile := os.Getenv("LOG_FILE")
	utils.LoggingSettings(logfile)
}

func main() {
	db := mysql.DB
	// TODO 要確認
	// defer func(db *sql.DB) {
	// 	err := db.Close()
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }(db)

	// DI
	articleRepo := rdb.NewArticleRepository(db)
	authorRepo := rdb.NewAuthorRepository(db)
	articleService := service.NewArticleService(articleRepo)
	articleUsecase := usecase.NewArticleUsecase(articleRepo, authorRepo, articleService, time.Second*120)
	err := handler2.NewHandler(articleUsecase)
	if err != nil {
		log.Fatalln(err)
	}
}
