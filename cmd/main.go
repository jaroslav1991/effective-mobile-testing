package main

import (
	_ "effective_mobile_testing/docs"
	"effective_mobile_testing/internal/config"
	"effective_mobile_testing/internal/connection"
	"effective_mobile_testing/internal/handlers"
	"effective_mobile_testing/internal/service"
	"effective_mobile_testing/internal/service/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"log/slog"
	"os"
)

// @title           Time-tracker API
// @version         1.0
// @description     Effective mobile testing
// @host			localhost:8080
// @BasePath		/

func main() {
	loggerSlog := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(loggerSlog)

	if err := godotenv.Load("./.env"); err != nil {
		slog.Error("Error loading .env file", slog.String("err", err.Error()))
	}

	dbConfig := config.GetDBConfig()
	slog.Debug(fmt.Sprintf("dbConfig: %v", dbConfig))

	db, err := connection.NewPostgresDB(dbConfig)
	if err != nil {
		return
	}
	slog.Debug("postgres connection created")

	repo := repository.NewUserTaskRepo(db)
	userTaskService := service.NewUserTaskService(repo)
	handler := handlers.NewHandlers(userTaskService)

	if err := connection.InitSchema(db); err != nil {
		return
	}
	slog.Debug("schema initialized")

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/user/create/", handler.CreateUser())
	r.PATCH("/user/start-tracking/", handler.StartTracking())
	r.PATCH("/user/stop-tracking/", handler.StopTracking())
	r.GET("/user/get-costs/", handler.GetLaborCosts())
	r.GET("/users/", handler.GetUserByFilters())
	r.DELETE("/user/", handler.DeleteUser())
	r.PATCH("/user/", handler.UpdateUser())

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
