package main

import (
	"log"

	"github.com/KCFLEX/Taxi-user-service/internal/config"
	"github.com/KCFLEX/Taxi-user-service/internal/handlers"
	"github.com/KCFLEX/Taxi-user-service/internal/repository"
	services "github.com/KCFLEX/Taxi-user-service/internal/service"

	"github.com/KCFLEX/Taxi-user-service/internal/service/tokens"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Print(err)
		return
	}
	repo, err := repository.New(config)
	if err != nil {
		log.Panic(err)
	}
	defer repo.Close()
	defer repo.CloseRedis()
	token := tokens.New(config)
	srv := services.New(repo, token)

	Handler := handlers.New(config, srv)
	Handler.Serve()

}
