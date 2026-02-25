package main

import (
	"context"
	"log"
	"net/http"

	"ApiGolangv2/internal/config"
	"ApiGolangv2/internal/handler"
	"ApiGolangv2/internal/repository"
	"ApiGolangv2/internal/routes"
	"ApiGolangv2/internal/service"

	"github.com/joho/godotenv"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on container env variables")
	}
	// start Datadog tracer pertama sebelum apapun
	tracer.Start(
		tracer.WithServiceName("task-manager"),
		tracer.WithEnv("development"),
		tracer.WithServiceVersion("1.0.0"),
	)
	defer tracer.Stop()

	ctx := context.Background()

	// init DB & Redis
	config.DatabaseConnection(ctx)
	redisClient, err := config.InitRedis()
	if err != nil {
		log.Fatal(err)
	}

	taskRepo := repository.NewTaskRepository(config.DB)
	taskService := service.NewTaskService(taskRepo, redisClient)
	taskHandler := handler.NewTaskHandler(taskService)

	// router dengan auto instrument
	router := routes.NewRouter(taskHandler)

	log.Println(" Server running on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
