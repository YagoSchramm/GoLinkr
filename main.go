package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	usecaseimpl "github.com/YagoSchramm/Golinkr/domain/usecase/impl"
	"github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository"
	repoimpl "github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository/impl"
	"github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository/util"
	approuter "github.com/YagoSchramm/Golinkr/infrastructure/router"
	"github.com/YagoSchramm/Golinkr/infrastructure/router/modules"
	"github.com/gorilla/mux"
)

func main() {
	router, cleanup, err := buildApp()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func buildApp() (*mux.Router, func(), error) {
	loadDotEnv()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, func() {}, errors.New("DATABASE_URL is not set")
	}

	db, err := util.NewPostgresConnection(dsn)
	if err != nil {
		return nil, func() {}, err
	}

	linkRepository := repository.LinkRepository(repoimpl.NewLinkRepository(db))
	cleanup := func() {
		_ = db.Close()
	}

	linkUsecase := usecaseimpl.NewLinkUsecase(linkRepository)
	linkModule := modules.NewLinkModule(linkUsecase)

	router := mux.NewRouter()
	approuter.Mount(router, linkModule)

	return router, cleanup, nil
}

func loadDotEnv() {
	content, err := os.ReadFile(".env")
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)

		if key == "" {
			continue
		}

		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
}
