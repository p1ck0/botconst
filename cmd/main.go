package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func main() {
	cfg, err := config.NewConfig(configPath)
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
