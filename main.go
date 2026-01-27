package main

import (
	"cv-landing/pkg/activity"
	"cv-landing/pkg/files"
	"cv-landing/pkg/handlers"
	"cv-landing/pkg/middleware"
	"cv-landing/pkg/tags"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func mustGetEnv(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		panic(fmt.Sprintf("the value of %s is empty", envName))
	}
	return env
}

func main() {
	apiVersion := "v1"

	user := mustGetEnv("POSTGRES_USER")
	password := mustGetEnv("POSTGRES_PASSWORD")
	host := mustGetEnv("POSTGRES_HOST")
	frontendDomain := mustGetEnv("FRONTEND_DOMAIN")

	dsn := fmt.Sprintf("user=%s password=%s host=%s sslmode=disable", user, password, host)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}

	activity := handlers.ActivityHandler{
		Repo: activity.ActivityHandler{
			DB: db,
		},
	}
	skills := handlers.SkillsHandler{
		Repo: files.FileHandler{
			BasePath: []string{"public"},
		},
	}
	tags := handlers.TagsHandler{
		Repo: tags.TagsHandler{
			DB: db,
		},
	}

	v1ApiRouter := mux.NewRouter().PathPrefix("/" + apiVersion).Subrouter()
	v1ApiRouter.HandleFunc("/activity/{type:[[:alpha:]]+}/", activity.Get).Methods("GET")
	v1ApiRouter.HandleFunc("/skills/{type:[[:alpha:]]+}/", skills.Get).Methods("GET")
	v1ApiRouter.HandleFunc("/tags/{type:[[:alpha:]]+}/", tags.Get).Methods("GET")
	v1ApiRouter.HandleFunc("/tags/{id:\\d+}/{type:[[:alpha:]]+}/", tags.Get).Methods("GET")

	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.EnableCors(v1ApiRouter, frontendDomain),
	}
	fmt.Println("Starting a server...")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
