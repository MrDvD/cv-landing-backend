package main

import (
	"cv-landing/pkg/activity"
	"cv-landing/pkg/files"
	"cv-landing/pkg/handlers"
	"database/sql"
	"fmt"
	"net/http"
	"os"

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

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/activity/{type}/", activity.Get)
	apiMux.HandleFunc("/skills/{type}/", skills.Get)

	apiHandler := http.StripPrefix("/"+apiVersion, apiMux)
	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("/%s/", apiVersion), apiHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Starting a server...")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
