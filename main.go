package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/maxoov1/faq-api/docs"
	"github.com/maxoov1/faq-api/pkg/auth"
	"github.com/maxoov1/faq-api/pkg/config"
	"github.com/maxoov1/faq-api/pkg/database/mongodb"
	"github.com/maxoov1/faq-api/pkg/handler"
	"github.com/maxoov1/faq-api/pkg/hash"
	"github.com/maxoov1/faq-api/pkg/repository"
	"github.com/maxoov1/faq-api/pkg/server"
	"github.com/maxoov1/faq-api/pkg/service"
)

const (
	timeout    = time.Second * 5
	configPath = "configs/config.yaml"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.NewConfig(configPath)
	config.Host = cfg.Host
	if err != nil {
		log.Fatal(err)
	}

	mongoClient, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.Username, cfg.Mongo.Password)
	if err != nil {
		log.Fatal(err)
	}

	db := mongoClient.Database(cfg.Mongo.Name)
	hasher := hash.NewMD5Hasher(cfg.Salt)

	manager, err := auth.NewManager(cfg.TokenManager.SigningKey)
	if err != nil {
		log.Fatal(err)
	}

	repos := repository.NewRepositories(db)
	services := service.NewServices(repos, hasher, manager, cfg.TokenManager.TokenTTL)
	handlers := handler.NewHandler(services)

	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Fatal(err)
	}

	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}
