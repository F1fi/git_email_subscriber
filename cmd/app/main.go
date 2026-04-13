package main

import (
	"fmt"
	"context"
	"time"
	"git_email_subscriber/internal/api"
	"git_email_subscriber/internal/github_api"
	"git_email_subscriber/internal/email_service"
	"git_email_subscriber/internal/subscriptions"
	"git_email_subscriber/internal/db"
	"git_email_subscriber/internal/scanner"
)

func main() {
	// fmt.Println("Start")

	r := gin.Default()

	dbPool := db.NewPostgresPool()
	db.RunMigrations(dbPool)

	repo := subscriptions.CreateRepo(dbPool)
	ghApi := github_api.CreateGitHubApi("token")
	emailService := email_service.CreateEmailService()
	subscriptionService := subscriptions.NewSubscriptionService(ghApi, repo, emailService)

	ctx := context.Background()

	scanner := scanner.NewScanner(
		repo,
		ghApi,
		emailService,
		5*time.Minute,
	)

	scanner.Start(ctx)

	handler := api.NewHandler(subscriptionService)
	fmt.Println(handler)

	handler.RegisterRoutes(r)

	r.Run(":8080")

	// fmt.Println("End")
}