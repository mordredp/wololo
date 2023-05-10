package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/mordredp/wololo/auth"
)

// Global variables
var appConfig AppConfig
var appData AppData

func main() {

	setWorkingDir()

	loadConfig()

	loadData()

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(auth.Identify)
	router.Use(auth.Clean)

	router.Get("/", renderHomePage)
	router.Get("/logout", auth.Logout)
	router.Post("/login", auth.Login)
	router.Get("/refresh", auth.Refresh)

	httpListen := appConfig.IP + ":" + strconv.Itoa(appConfig.Port)
	log.Printf("Startup Webserver on \"%s\"", httpListen)

	log.Fatal(http.ListenAndServe(httpListen, router))
}

func setWorkingDir() {

	thisApp, err := os.Executable()
	if err != nil {
		log.Fatalf("Error determining the directory. \"%s\"", err)
	}
	appPath := filepath.Dir(thisApp)
	os.Chdir(appPath)
	log.Printf("Set working directory: %s", appPath)

}

func loadConfig() {

	err := cleanenv.ReadConfig("config.json", &appConfig)
	if err != nil {
		log.Fatalf("Error loading config.json file. \"%s\"", err)
	}
	log.Printf("Application configuratrion loaded from config.json")

}
