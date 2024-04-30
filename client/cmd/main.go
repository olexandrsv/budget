package main

import (
	"budget/pkg/repository"
	"budget/pkg/service"
	"budget/pkg/ui"
)

func main() {
	cache := repository.NewCache()
	db := repository.NewDatabase()
	repo := repository.NewRepository(db, cache)
	srv := service.New(repo)
	u := ui.New("../assets/leds", srv)
	u.Run()
}