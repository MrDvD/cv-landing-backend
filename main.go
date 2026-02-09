package main

import (
	"cv-landing-backend/pkg/activity"
	"cv-landing-backend/pkg/attachments"
	"cv-landing-backend/pkg/files"
	"cv-landing-backend/pkg/handlers"
	"cv-landing-backend/pkg/middleware"
	"cv-landing-backend/pkg/tags"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync"

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
	publicDomain := mustGetEnv("PUBLIC_DOMAIN")
	privateDomain := mustGetEnv("PRIVATE_DOMAIN")

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
			Path: []string{"public"},
		},
	}
	tags := handlers.TagsHandler{
		Repo: tags.TagsHandler{
			DB: db,
		},
	}
	attachments := handlers.AttachmentHandler{
		Repo: attachments.AttachmentHandler{
			DB: db,
		},
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		publicRouter := mux.NewRouter().PathPrefix("/" + apiVersion).Subrouter()
		publicRouter.HandleFunc("/activity/{type:[[:alpha:]]+}/", activity.Get).Methods("GET")
		publicRouter.HandleFunc("/skills/{type:[[:alpha:]]+}/", skills.Get).Methods("GET")
		publicRouter.HandleFunc("/tags/{type:[[:alpha:]]+}/", tags.Get).Methods("GET")
		publicRouter.HandleFunc("/tags/{id:\\d+}/{type:[[:alpha:]]+}/", tags.Get).Methods("GET")
		publicRouter.HandleFunc("/attachments/{id:\\d+}/", attachments.Get).Methods("GET")

		publicServer := http.Server{
			Addr:    ":8080",
			Handler: middleware.EnableCors(publicRouter, publicDomain),
		}
		fmt.Println("starting a public server...")
		err = publicServer.ListenAndServe()
		if err != nil {
			fmt.Println("public server error:", err)
			return
		}
	}()
	go func() {
		defer wg.Done()
		privateRouter := mux.NewRouter().PathPrefix("/" + apiVersion).Subrouter()
		privateRouter.HandleFunc("/activity/", activity.Add).Methods("POST")
		privateRouter.HandleFunc("/tags/", tags.Add).Methods("POST")
		privateRouter.HandleFunc("/attachments/", attachments.Add).Methods("POST")

		privateRouter.HandleFunc("/activity/{id:\\d+}/", activity.Remove).Methods("DELETE")
		privateRouter.HandleFunc("/tags/{id:\\d+}/", tags.Remove).Methods("DELETE")
		privateRouter.HandleFunc("/attachments/{id:\\d+}/", attachments.Remove).Methods("DELETE")

		privateRouter.HandleFunc("/activity/{id:\\d+}/", activity.Edit).Methods("PATCH")
		privateRouter.HandleFunc("/tags/{id:\\d+}/", tags.Edit).Methods("PATCH")
		privateRouter.HandleFunc("/attachments/{id:\\d+}/", attachments.Edit).Methods("PATCH")

		privateServer := http.Server{
			Addr:    privateDomain,
			Handler: privateRouter,
		}
		fmt.Println("starting a private server...")
		err = privateServer.ListenAndServe()
		if err != nil {
			fmt.Println("private server error:", err)
			return
		}
	}()
	wg.Wait()
}
